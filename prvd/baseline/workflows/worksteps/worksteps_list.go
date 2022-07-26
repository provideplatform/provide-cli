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

package worksteps

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var listBaselineWorkstepsCmd = &cobra.Command{
	Use:   "list",
	Short: "List baseline worksteps",
	Long:  `List all available baseline workflow worksteps`,
	Run:   listWorksteps,
}

func listWorksteps(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listWorkstepsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		fmt.Printf("failed to list worksteps; %s", err.Error())
		os.Exit(1)
	}

	if workflowID == "" {
		workflowPrompt(*token.AccessToken)
	}

	worksteps, err := baseline.ListWorksteps(*token.AccessToken, workflowID, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to list worksteps; %s", err.Error())
		os.Exit(1)
	}

	if len(worksteps) == 0 {
		fmt.Print("No worksteps found\n")
		return
	}

	for _, workstep := range worksteps {
		result, _ := json.MarshalIndent(workstep, "", "\t")
		fmt.Printf("%s\n", string(result))
	}
}

func init() {
	listBaselineWorkstepsCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	listBaselineWorkstepsCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	listBaselineWorkstepsCmd.Flags().StringVar(&workflowID, "workflow", "", "workflow identifier")

	listBaselineWorkstepsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineWorkstepsCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of baseline workgroups to retrieve per page")
	listBaselineWorkstepsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
