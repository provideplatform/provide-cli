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
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/axiom"

	"github.com/spf13/cobra"
)

var workflowID string

var detailBaselineWorkflowCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific axiom workflow",
	Long:  `Retrieve details for a specific axiom workflow by identifier, scoped to the authorized API token`,
	Run:   fetchWorkflowDetails,
}

func fetchWorkflowDetails(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepDetails)
}

func fetchWorkflowDetailsRun(cmd *cobra.Command, args []string) {
	if err := common.RequireOrganization(); err != nil {
		fmt.Printf("failed to retrive workflow details; %s", err.Error())
		os.Exit(1)
	}

	if err := common.RequireWorkgroup(); err != nil {
		fmt.Printf("failed to retrive workflow details; %s", err.Error())
		os.Exit(1)
	}

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to retrieve workflow details; %s", err.Error())
		os.Exit(1)
	}

	if workflowID == "" {
		workflowPrompt(*token.AccessToken)
	}

	w, err := axiom.GetWorkflowDetails(*token.AccessToken, workflowID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to retrieve workflow details; %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(w, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func workflowPrompt(token string) {
	workflows, err := axiom.ListWorkflows(token, map[string]interface{}{
		"workgroup_id": common.WorkgroupID,
	})
	if err != nil {
		fmt.Printf("failed to retrieve workflow details; %s", err.Error())
		os.Exit(1)
	}

	if len(workflows) == 0 {
		fmt.Print("No workflows found\n")
		os.Exit(1)
	}

	opts := make([]string, 0)

	for _, workflow := range workflows {
		opts = append(opts, *workflow.Name)
	}

	prompt := promptui.Select{
		Label: "Select Workflow",
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("failed to retrieve workflow details; %s", err.Error())
		os.Exit(1)
	}

	workflowID = workflows[i].ID.String()
}

func init() {
	detailBaselineWorkflowCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	detailBaselineWorkflowCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	detailBaselineWorkflowCmd.Flags().StringVar(&workflowID, "workflow", "", "workflow identifier")
}
