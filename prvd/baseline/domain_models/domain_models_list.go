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
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var name string // FIXME-- using var 'name' instead of 'type' bc type is reserved keyword

var page uint64
var rpp uint64

var listBaselineDomainModelsCmd = &cobra.Command{
	Use:   "list",
	Short: "List baseline domain models",
	Long:  `List all available baseline domain models`,
	Run:   listDomainModels,
}

func listDomainModels(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listDomainModelsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}

	if common.WorkgroupID == "" {
		prompt := promptui.Prompt{
			IsConfirm: true,
			Label:     "Select workgroup",
		}

		_, err := prompt.Run()
		if err == nil {
			common.RequireWorkgroup()
		}
	}

	var ref string
	if name == "" {
		prompt := promptui.Prompt{
			IsConfirm: true,
			Label:     "Enter model type",
		}

		_, err := prompt.Run()
		if err == nil {
			listModelsTypePrompt()

			ref = common.SHA256(fmt.Sprintf("%s.%s", common.OrganizationID, name))
		}

	}

	common.AuthorizeOrganizationContext(true)

	token := common.RequireOrganizationToken()

	models, err := baseline.ListMappings(token, map[string]interface{}{
		"workgroup_id": common.WorkgroupID,
		"ref":          ref,
		"page":         fmt.Sprintf("%d", page),
		"rpp":          fmt.Sprintf("%d", rpp),
	})
	if err != nil {
		log.Printf("failed to retrieve baseline domain models; %s", err.Error())
		os.Exit(1)
	}

	if len(models) == 0 {
		fmt.Print("No domain models found\n")
		return
	}

	for _, model := range models {
		result := fmt.Sprintf("%s\t%d field(s)\t%s\n", model.ID.String(), len(model.Models[0].Fields), *model.Type)
		fmt.Print(result)
	}
}

func listModelsTypePrompt() {
	prompt := promptui.Prompt{
		Label: "Model Type",
	}

	result, err := prompt.Run()
	if err == nil {
		name = result
	}
}

func init() {
	listBaselineDomainModelsCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	listBaselineDomainModelsCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	listBaselineDomainModelsCmd.Flags().StringVar(&name, "type", "", "domain model type")
	listBaselineDomainModelsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineDomainModelsCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of baseline domain models to retrieve per page")
	listBaselineDomainModelsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
