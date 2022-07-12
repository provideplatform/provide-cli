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

package workgroups

import (
	"encoding/json"
	"fmt"
	"os"

	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-cli/prvd/organizations"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var updateBaselineWorkgroupCmd = &cobra.Command{
	Use:   "update",
	Short: "Update baseline workgroup",
	Long:  `Update baseline workgroup`,
	Run:   updateWorkgroup,
}

func updateWorkgroup(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepUpdate)
}

func updateWorkgroupRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	var localWg Workgroup
	raw, _ := json.Marshal(common.Workgroup)
	json.Unmarshal(raw, &localWg)

	var localOrg organizations.Organization
	raw, _ = json.Marshal(common.Organization)
	json.Unmarshal(raw, &localOrg)

	wgID, err := uuid.FromString(common.WorkgroupID)
	if err != nil {
		fmt.Printf("failed to update baseline workgroup: %s", err.Error())
		os.Exit(1)
	}

	isOperator := localOrg.Metadata.Workgroups[wgID].OperatorSeparationDegree == 0

	if err := updateWorkgroupPrompt(&localWg, isOperator); err != nil {
		fmt.Printf("failed to update baseline workgroup: %s", err.Error())
		os.Exit(1)
	}

	var wgParams map[string]interface{}
	raw, _ = json.Marshal(localWg)
	json.Unmarshal(raw, &wgParams)

	token := common.RequireOrganizationToken()

	if err := baseline.UpdateWorkgroup(token, common.WorkgroupID, wgParams); err != nil {
		fmt.Printf("failed to update baseline workgroup: %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(wgParams, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func updateWorkgroupPrompt(wg *Workgroup, isOperator bool) error {
	// name
	if name == "" {
		prompt := promptui.Prompt{
			Label:   "Workgroup Name",
			Default: *wg.Name,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("name cannot be empty")
				}

				return nil
			},
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
		}

		name = result
	}
	*wg.Name = name

	// description
	if description == "" {
		var defaultDesc string
		if wg.Description != nil {
			defaultDesc = *wg.Description
		}

		prompt := promptui.Prompt{
			Label:   "Workgroup Description",
			Default: defaultDesc,
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
		}

		description = result
	}
	*wg.Description = description

	// TODO-- vault, systems

	// layers
	if isOperator {
		if common.NetworkID == "" {
			common.RequireL1Network()
		}
		uuidNetworkID, err := uuid.FromString(common.NetworkID)
		if err != nil {
			fmt.Printf("failed to update baseline workgroup: %s", err.Error())
			os.Exit(1)
		}
		*wg.NetworkID = uuidNetworkID

		if common.L2NetworkID == "" {
			common.RequireL2Network()
		}
		uuidL2NetworkID, err := uuid.FromString(common.L2NetworkID)
		if err != nil {
			fmt.Printf("failed to update baseline workgroup: %s", err.Error())
			os.Exit(1)
		}
		*wg.Config.L2NetworkID = uuidL2NetworkID
	} else if !isOperator && (common.NetworkID != "" || common.L2NetworkID != "") {
		return fmt.Errorf("workgroup participants cannot update layers\n")
	}

	return nil
}

func init() {
	updateBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	updateBaselineWorkgroupCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	updateBaselineWorkgroupCmd.Flags().StringVar(&name, "name", "", "name of the baseline workgroup")
	updateBaselineWorkgroupCmd.Flags().StringVar(&description, "description", "", "description of the baseline workgroup")
	updateBaselineWorkgroupCmd.Flags().StringVar(&common.NetworkID, "network", "", "nchain network id of the baseline mainnet to use for this workgroup")
	updateBaselineWorkgroupCmd.Flags().StringVar(&common.L2NetworkID, "l2", "", "nchain l2 network id of the baseline layer 2 to use for this workgroup")
}
