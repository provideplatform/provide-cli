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

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var deployBaselineWorkflowCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy baseline workflow",
	Long:  `deploy a baseline prototype workflow`,
	Run:   deployWorkflow,
}

func deployWorkflow(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepDeploy)
}

func deployWorkflowRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to deploy workflow; %s", err.Error())
		os.Exit(1)
	}
	if workflowID == "" {
		workflowPrompt(*token.AccessToken)
	}

	w, err := baseline.GetWorkflowDetails(*token.AccessToken, workflowID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to deploy workflow; %s", err.Error())
		os.Exit(1)
	}
	if w.WorkflowID != nil {
		log.Print("failed to deploy workflow; cannot deploy a workflow instance")
		os.Exit(1)
	}
	if *w.Status != "draft" {
		log.Print("failed to deploy workflow; cannot deploy a non-draft instance")
		os.Exit(1)
	}

	ws, err := baseline.ListWorksteps(*token.AccessToken, workflowID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to deploy workflow; %s", err.Error())
		os.Exit(1)
	}

	hasFinality := false
	for _, workstep := range ws {
		var metadata map[string]interface{}
		raw, _ := json.Marshal(workstep.Metadata)
		json.Unmarshal(raw, &metadata)

		if metadata["prover"] == nil {
			log.Printf("failed to deploy workflow; all worksteps must have a prover")
			os.Exit(1)
		}

		if workstep.RequireFinality {
			hasFinality = true
		}
	}

	if !hasFinality {
		log.Printf("failed to deploy workflow; at least 1 workstep must require finality")
		os.Exit(1)
	}

	deployed, err := baseline.DeployWorkflow(*token.AccessToken, workflowID, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to deploy workflow; %s", err.Error())
		os.Exit(1)
	}

	// wait til status is deployed ?

	result, _ := json.MarshalIndent(deployed, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func init() {
	deployBaselineWorkflowCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	deployBaselineWorkflowCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	deployBaselineWorkflowCmd.Flags().StringVar(&workflowID, "workflow", "", "workflow identifier")

	deployBaselineWorkflowCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
