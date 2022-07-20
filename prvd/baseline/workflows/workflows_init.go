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
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var name string
var description string
var version string

var Optional bool

var initBaselineWorkflowCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize baseline workflow",
	Long:  `Initialize and configure a new baseline workflow`,
	Run:   initWorkflow,
}

func initWorkflow(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func initWorkflowRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
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
			fmt.Printf("failed to initialize workflow; %s", err.Error())
			os.Exit(1)
		}
	}

	common.AuthorizeOrganizationContext(true)

	token := common.RequireOrganizationToken()

	params := map[string]interface{}{
		"workgroup_id": common.WorkgroupID,
		"name":         name,
		"version":      version,
	}

	if description != "" {
		params["description"] = description
	}

	w, err := baseline.CreateWorkflow(token, params)
	if err != nil {
		fmt.Printf("failed to initialize workflow; %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(w, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func namePrompt() {
	prompt := promptui.Prompt{
		Label: "Workflow Name",
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
		Label: "Workflow Description",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	description = result
}

func versionPrompt() {
	prompt := promptui.Prompt{
		Label:   "Workflow Version",
		Default: "0.0.1",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("version is required")
			}

			if _, err := semver.Make(s); err != nil {
				return err
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	version = result
}

func init() {
	initBaselineWorkflowCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineWorkflowCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")

	initBaselineWorkflowCmd.Flags().StringVar(&name, "name", "", "name of the baseline workflow")
	initBaselineWorkflowCmd.Flags().StringVar(&description, "description", "", "description of the baseline workflow")
	initBaselineWorkflowCmd.Flags().StringVar(&version, "version", "", "version of the baseline workflow")

	initBaselineWorkflowCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
