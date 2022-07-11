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

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var description string

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
	if name == "" {
		createModelTypePrompt()
	}
	if description == "" {
		descriptionPrompt()
	}

	fields := make([]*baseline.MappingField, 0)

	fieldsPrompt(&fields)

	pk, err := primaryKeyPrompt(fields)
	if err != nil {
		log.Printf("failed to initialize baseline domain model; %s", err.Error())
		os.Exit(1)
	}

	common.AuthorizeOrganizationContext(true)

	token := common.RequireOrganizationToken()

	m, err := baseline.CreateMapping(token, map[string]interface{}{
		"name":        name,
		"description": description,
		"type":        name,
		"models": []interface{}{
			map[string]interface{}{
				"type":        name,
				"description": description,
				"fields":      fields,
				"primary_key": pk,
			},
		},
		"workgroup_id": common.WorkgroupID,
	})
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

func primaryKeyPrompt(fields []*baseline.MappingField) (string, error) {
	fieldNames := make([]string, 0)
	for _, field := range fields {
		fieldNames = append(fieldNames, field.Name)
	}

	prompt := promptui.Select{
		Label: "Select Primary Key",
		Items: fieldNames, // TODO-- use templates
	}

	i, _, err := prompt.Run()
	if err == nil {
		for j, field := range fields {
			if i == j {
				field.IsPrimaryKey = true
				return field.Name, nil
			}
		}
	}

	return "", err
}

func init() {
	initBaselineDomainModelCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineDomainModelCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	initBaselineDomainModelCmd.Flags().StringVar(&name, "type", "", "model type")
	initBaselineDomainModelCmd.Flags().StringVar(&description, "description", "", "model description")

	initBaselineDomainModelCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
