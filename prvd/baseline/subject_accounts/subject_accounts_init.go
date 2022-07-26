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

package subject_accounts

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/spf13/cobra"
)

var Optional bool
var paginate bool

var initBaselineSubjectAccountCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize baseline subject account",
	Long:  `Initialize and configure a new baseline subject account`,
	Run:   createSubjectAccount,
}

func createSubjectAccount(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func createSubjectAccountRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}

	// TODO-- check if user can pass workgroup id of workgroup that is not associated with organization id and handle that
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	if common.NetworkID == "" {
		common.RequireL1Network()
	}
	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()

	sa, err := baseline.CreateSubjectAccount(*token.AccessToken, common.OrganizationID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"organization_id":            common.OrganizationID,
			"organization_address":       common.Organization.Metadata.Address,
			"organization_refresh_token": *token.RefreshToken,
			"workgroup_id":               common.WorkgroupID,
			"registry_contract_address":  common.Organization.Metadata.Address,
			"network_id":                 common.NetworkID,
		},
	})
	if err != nil {
		log.Printf("failed to initialize baseline subject account; %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(sa, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func init() {
	initBaselineSubjectAccountCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineSubjectAccountCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	initBaselineSubjectAccountCmd.Flags().StringVar(&common.NetworkID, "network", "", "nchain network id of the baseline mainnet to use for this workgroup")

	initBaselineSubjectAccountCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
