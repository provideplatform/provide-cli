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

package workflows

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var filterPrototypes bool
var filterInstances bool

var page uint64
var rpp uint64

var listBaselineWorkflowsCmd = &cobra.Command{
	Use:   "list",
	Short: "List baseline workflows",
	Long:  `List all available baseline workflows`,
	Run:   listWorkflows,
}

func listWorkflows(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listWorkflowsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}
	if !filterInstances {
		prompt := promptui.Prompt{
			IsConfirm: true,
			Label:     "Filter Instances",
		}

		if _, err := prompt.Run(); err == nil {
			filterInstances = true
		}
	}
	if !filterPrototypes {
		prompt := promptui.Prompt{
			IsConfirm: true,
			Label:     "Filter Prototypes",
		}

		if _, err := prompt.Run(); err == nil {
			filterPrototypes = true
		}
	}

	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		fmt.Printf("failed to list workflows; %s", err.Error())
		os.Exit(1)
	}

	params := map[string]interface{}{
		"workgroup_id": common.WorkgroupID,
	}

	if filterInstances {
		params["filter_instances"] = "true"
	}

	if filterPrototypes {
		params["filter_prototypes"] = "true"
	}

	workflows, err := baseline.ListWorkflows(*token.AccessToken, params)
	if err != nil {
		fmt.Printf("failed to list workflows; %s", err.Error())
		os.Exit(1)
	}

	if len(workflows) == 0 {
		fmt.Print("No workflows found\n")
		return
	}

	for _, workflow := range workflows {
		result, _ := json.MarshalIndent(workflow, "", "\t")
		fmt.Printf("%s\n", string(result))
	}
}

func init() {
	listBaselineWorkflowsCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	listBaselineWorkflowsCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")

	listBaselineWorkflowsCmd.Flags().BoolVar(&filterInstances, "filter-instances", false, "filter workflow prototypes")
	listBaselineWorkflowsCmd.Flags().BoolVar(&filterPrototypes, "filter-prototypes", false, "filter workflow instances")

	listBaselineWorkflowsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineWorkflowsCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of baseline workgroups to retrieve per page")
	listBaselineWorkflowsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
