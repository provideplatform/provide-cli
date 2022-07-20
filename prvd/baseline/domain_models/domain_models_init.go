/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package domain_models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var isSchema bool
var schemaQuery string

var description string

var fields string
var primaryKey string

var Optional bool
var paginate bool

var initBaselineDomainModelCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize baseline domain model",
	Long:  `Initialize and configure a new baseline domain model`,
	Run:   initDomainModel,
}

func initDomainModel(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func initDomainModelRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to initialize baseline domain model; %s", err.Error())
		os.Exit(1)
	}

	hasSystems := len(common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs) > 0

	if hasSystems && !isSchema {
		isSchemaPrompt()
	} else if hasSystems && isSchema {
		fmt.Print("failed to initialize baseline domain model; cannot create a domain model from a schema without systems")
		os.Exit(1)
	}

	var params map[string]interface{}
	if isSchema {
		schemaQueryPrompt()

		vaultID := common.Organization.Metadata.Workgroups[common.Workgroup.ID].VaultID
		systemIDs := common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs

		isOperator := common.Organization.Metadata.Workgroups[common.Workgroup.ID].OperatorSeparationDegree == 0
		if isOperator {
			vaultID = common.Workgroup.Config.VaultID
			systemIDs = common.Workgroup.Config.SystemSecretIDs
		}

		IDs := make([]string, 0)
		for _, ID := range systemIDs {
			IDs = append(IDs, ID.String())
		}

		schemas, err := baseline.ListSchemas(*token.AccessToken, common.WorkgroupID, map[string]interface{}{
			"vault_id":          vaultID.String(),
			"system_secret_ids": strings.Join(IDs, ","),
			"q":                 schemaQuery,
		})
		if err != nil {
			log.Printf("failed to initialize baseline domain model; %s", err.Error())
			os.Exit(1)
		}

		schemaOpts := make([]string, 0)
		for _, schema := range schemas {
			schemaOpts = append(schemaOpts, *schema.Name)
		}

		prompt := promptui.Select{
			Label: "Select Schema",
			Items: schemaOpts,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize baseline domain model; %s", err.Error())
			os.Exit(1)
		}

		ref := common.SHA256(fmt.Sprintf("%s.%s", common.OrganizationID, schemaOpts[i]))
		models, err := baseline.ListMappings(*token.AccessToken, map[string]interface{}{
			"workgroup_id": common.WorkgroupID,
			"ref":          ref,
			// "page":         fmt.Sprintf("%d", page),
			// "rpp":          fmt.Sprintf("%d", rpp),
		})

		if len(models) > 0 {
			fmt.Print("failed to initialize baseline domain model; schema mapping exists")
			os.Exit(1)
		}

		schema, err := baseline.GetSchemaDetails(*token.AccessToken, common.OrganizationID, ref, map[string]interface{}{})
		if err != nil {
			fmt.Printf("failed to initialize baseline domain model; %s", err.Error())
			os.Exit(1)
		}

		fields := make([]interface{}, 0)
		for _, field := range schema.Fields {
			var f map[string]interface{}
			raw, _ := json.Marshal(field)
			json.Unmarshal(raw, &f)

			f["type"] = "string"
			fields = append(fields, f)
		}

		model := map[string]interface{}{
			"type":        *schema.Type,
			"fields":      fields,
			"primary_key": "",
			"standard":    "sap",
		}

		params = map[string]interface{}{
			"name":         schema.Name,
			"description":  schema.Description,
			"type":         schema.Type,
			"models":       []interface{}{model},
			"workgroup_id": common.WorkgroupID,
		}
	} else {
		if name == "" {
			createModelTypePrompt()
		}
		if description == "" {
			descriptionPrompt()
		}

		localFields := make([]*baseline.MappingField, 0)
		if fields != "" {
			if err := json.Unmarshal([]byte(fields), &localFields); err != nil {
				log.Printf("failed to initialize baseline domain model; %s", err.Error())
				os.Exit(1)
			}

			if err := validateFields(localFields); err != nil {
				log.Printf("failed to initialize baseline domain model; %s", err.Error())
				os.Exit(1)
			}
		}

		fieldsPrompt(&localFields)

		if err := primaryKeyPrompt(localFields); err != nil {
			log.Printf("failed to initialize baseline domain model; %s", err.Error())
			os.Exit(1)
		}

		modelParam := map[string]interface{}{
			"type":        name,
			"fields":      localFields,
			"primary_key": primaryKey,
		}

		if description != "" {
			modelParam["description"] = description
		}

		params = map[string]interface{}{
			"name": name,
			"type": name,
			"models": []interface{}{
				modelParam,
			},
			"workgroup_id": common.WorkgroupID,
		}
	}

	m, err := baseline.CreateMapping(*token.AccessToken, params)
	if err != nil {
		log.Printf("failed to initialize baseline domain model; %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(m, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func createModelTypePrompt() {
	prompt := promptui.Prompt{
		Label: "Model type",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("type cannot be empty")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	name = result
}

func descriptionPrompt() {
	prompt := promptui.Prompt{
		Label: "Model Description",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	description = result
}

func fieldsPrompt(fields *[]*baseline.MappingField) func([]*baseline.MappingField) {
	if len(*fields) > 0 {
		prompt := promptui.Prompt{
			IsConfirm: true,
			Label:     "Add Field",
		}

		_, err := prompt.Run()
		if err != nil {
			return nil
		}
	}

	prompt := promptui.Prompt{
		Label: "Field Name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("field name is required")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("failed to initialize baseline domain model; %s", err.Error())
		os.Exit(1)
	}

	// FIXME-- can probably do this more simply - export as const from provide-go ??
	fieldTypes := make([]string, 2)
	fieldTypes[0] = "Number"
	fieldTypes[1] = "String"

	selectPrompt := promptui.Select{
		Label: "Field Type",
		Items: fieldTypes,
	}

	i, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("failed to initialize baseline domain model; %s", err.Error())
		os.Exit(1)
	}

	*fields = append(*fields, &baseline.MappingField{
		Name: result,
		Type: fieldTypes[i],
	})
	return fieldsPrompt(fields)

}

func primaryKeyPrompt(fields []*baseline.MappingField) error {
	if primaryKey == "" {
		fieldNames := make([]string, 0)
		for _, field := range fields {
			fieldNames = append(fieldNames, field.Name)
		}

		prompt := promptui.Select{
			Label: "Select Primary Key",
			Items: fieldNames, // TODO-- use templates
		}

		i, _, err := prompt.Run()
		if err != nil {
			return err
		}

		primaryKey = fieldNames[i]
	}

	for _, field := range fields {
		if field.Name == primaryKey {
			field.IsPrimaryKey = true
			return nil
		}
	}

	return fmt.Errorf("primary key not found")
}

func validateFields(fields []*baseline.MappingField) error {
	for _, field := range fields {
		if field.Name == "" {
			return fmt.Errorf("field must have a name")
		}

		if field.Type != "number" && field.Type != "string" {
			return fmt.Errorf("%s is not a valid field type; fields must have type number or string", field.Type)
		}

		if field.IsPrimaryKey {
			return fmt.Errorf("cannot set primary key from the --fields flag; use the --primary-key flag instead") // TODO-- this should be supported
		}
	}

	return nil
}

func isSchemaPrompt() {
	prompt := promptui.Prompt{
		IsConfirm: true,
		Label:     "Create model from schema",
	}

	if _, err := prompt.Run(); err == nil {
		isSchema = true
	}
}

func schemaQueryPrompt() {
	if schemaQuery == "" {
		prompt := promptui.Prompt{
			Label:    "Schema Query",
			Validate: common.MandatoryValidation,
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize baseline domain model; %s", err.Error())
			os.Exit(1)
		}

		schemaQuery = result
	}
}

func init() {
	initBaselineDomainModelCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineDomainModelCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")

	initBaselineDomainModelCmd.Flags().BoolVar(&isSchema, "schema", false, "create a domain model from a schema")
	initBaselineDomainModelCmd.Flags().StringVar(&schemaQuery, "schema-query", "", "schema query string")

	initBaselineDomainModelCmd.Flags().StringVar(&name, "type", "", "model type")
	initBaselineDomainModelCmd.Flags().StringVar(&description, "description", "", "model description")
	initBaselineDomainModelCmd.Flags().StringVar(&fields, "fields", "", "model fields in the '[{\"name\": \"yourmother\", \"type\": \"string\"}, ...]' format")
	initBaselineDomainModelCmd.Flags().StringVar(&primaryKey, "primary-key", "", "model primary key")

	initBaselineDomainModelCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
