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

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var workflowID string

var name string
var description string
var requireFinality bool

var prover string

// TODO
// var modelID string
// var participants []string

var Optional bool

var initBaselineWorkstepCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize baseline workstep",
	Long:  `Initialize and configure a new baseline workstep`,
	Run:   initWorkstep,
}

func initWorkstep(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func initWorkstepRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		fmt.Printf("failed to initialize workstep; %s", err.Error())
		os.Exit(1)
	}

	if workflowID == "" {
		workflowPrompt(*token.AccessToken)
	}

	if name == "" {
		namePrompt()
	}
	if description == "" {
		descriptionPrompt()
	}

	provers := []map[string]interface{}{
		{
			"identifier":     "cubic",
			"name":           "General Consistency",
			"provider":       "gnark",
			"proving_scheme": "groth16",
			"curve":          "BN254",
		},
	}

	if prover == "" {
		proverPrompt()
	} else {
		isValid := false
		for _, p := range provers {
			if prover == p["name"] {
				isValid = true
			}
		}

		if !isValid {
			fmt.Print("failed to initialize workstep; invalid prover identifier")
			os.Exit(1)
		}
	}

	if !requireFinality {
		requireFinalityPrompt()
	}

	params := map[string]interface{}{
		"name":             name,
		"status":           "draft",
		"require_finality": requireFinality,
		"metadata": map[string]interface{}{
			"prover": provers[0],
		},
	}

	if description != "" {
		params["description"] = description
	}

	ws, err := baseline.CreateWorkstep(*token.AccessToken, workflowID, params)
	if err != nil {
		fmt.Printf("failed to initialize workstep; %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(ws, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func workflowPrompt(token string) {
	workflows, err := baseline.ListWorkflows(token, map[string]interface{}{
		"workgroup_id": common.WorkgroupID,
	})

	if len(workflows) == 0 {
		fmt.Print("No workflows found\n")
		os.Exit(1)
	}

	opts := make([]string, 0)
	for _, workflow := range workflows {
		opts = append(opts, *workflow.Name)
	}

	prompt := promptui.Select{
		Label: "Workflow",
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		os.Exit(1)
	}

	workflowID = workflows[i].ID.String()
}

func namePrompt() {
	prompt := promptui.Prompt{
		Label: "Workstep Name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("name is required")
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
		Label: "Workstep Description",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	description = result
}

func requireFinalityPrompt() {
	prompt := promptui.Prompt{
		Label:     "Require Finality",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err == nil {
		requireFinality = true
	}
}

func proverPrompt() {
	opts := []string{"General Consistency"}

	prompt := promptui.Select{
		Label: "Prover",
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		os.Exit(1)
	}

	prover = opts[i]
}

func init() {
	initBaselineWorkstepCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineWorkstepCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	initBaselineWorkstepCmd.Flags().StringVar(&workflowID, "workflow", "", "workflow identifier")

	initBaselineWorkstepCmd.Flags().StringVar(&name, "name", "", "name of the baseline workstep")
	initBaselineWorkstepCmd.Flags().StringVar(&description, "description", "", "description of the baseline workstep")
	initBaselineWorkstepCmd.Flags().BoolVar(&requireFinality, "require-finality", false, "require finality on the baseline workstep")

	initBaselineWorkstepCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
