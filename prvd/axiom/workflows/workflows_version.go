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

	"github.com/blang/semver/v4"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/axiom"
	"github.com/spf13/cobra"
)

var versionBaselineWorkflowCmd = &cobra.Command{
	Use:   "version",
	Short: "Version a axiom workflow",
	Long:  `Version an existing axiom workflow`,
	Run:   versionWorkflow,
}

func versionWorkflow(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepVersion)
}

func versionWorkflowRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		fmt.Printf("failed to version workflow; %s", err.Error())
		os.Exit(1)
	}
	if workflowID == "" {
		workflowPrompt(*token.AccessToken)
	}
	workflow, err := axiom.GetWorkflowDetails(*token.AccessToken, workflowID, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to version workflow; %s", err.Error())
		os.Exit(1)
	}
	if workflow.WorkflowID != nil {
		fmt.Print("failed to version workflow; cannot version a workflow instance")
		os.Exit(1)
	}
	if *workflow.Status != "deployed" {
		fmt.Print("failed to version workflow; cannot version a non-deployed workflow")
		os.Exit(1)
	}

	if name == "" {
		namePrompt()
	}
	if description == "" {
		descriptionPrompt()
	}
	if version == "" {
		versionPrompt()
	} else {
		if _, err := semver.Make(version); err != nil {
			fmt.Printf("failed to version workflow; %s", err.Error())
			os.Exit(1)
		}
	}

	v1, _ := semver.Make(*workflow.Version)
	v2, _ := semver.Make(version)

	if !v2.GT(v1) {
		fmt.Printf("failed to version workflow; new version must be greater than previous")
		os.Exit(1)
	}

	params := map[string]interface{}{
		"name":    name,
		"version": version,
	}

	if description != "" {
		params["description"] = description
	}

	w, err := axiom.VersionWorkflow(*token.AccessToken, workflowID, params)
	if err != nil {
		fmt.Printf("failed to version workflow; %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(w, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func init() {
	versionBaselineWorkflowCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	versionBaselineWorkflowCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	versionBaselineWorkflowCmd.Flags().StringVar(&workflowID, "workflow", "", "workflow identifier")

	versionBaselineWorkflowCmd.Flags().StringVar(&name, "name", "", "name of the axiom workflow")
	versionBaselineWorkflowCmd.Flags().StringVar(&description, "description", "", "description of the axiom workflow")
	versionBaselineWorkflowCmd.Flags().StringVar(&version, "version", "", "version of the axiom workflow")

	versionBaselineWorkflowCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
