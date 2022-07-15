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

package systems

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var listBaselineSystemsCmd = &cobra.Command{
	Use:   "list",
	Short: "List baseline systems",
	Long:  `List all available baseline systems of record`,
	Run:   listSystems,
}

func listSystems(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listSystemsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	vaultID := common.Organization.Metadata.Workgroups[common.Workgroup.ID].VaultID
	systemIDs := common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs

	isOperator := common.Organization.Metadata.Workgroups[common.Workgroup.ID].OperatorSeparationDegree == 0
	if isOperator {
		vaultID = common.Workgroup.Config.VaultID
		systemIDs = common.Workgroup.Config.SystemSecretIDs
	}

	common.AuthorizeOrganizationContext(true)

	token := common.RequireOrganizationToken()

	secrets := make([]*vault.Secret, 0)

	for _, secretID := range systemIDs {
		secret, err := vault.FetchSecret(token, vaultID.String(), secretID.String(), map[string]interface{}{})
		if err != nil {
			log.Printf("failed to retrieve systems; %s", err.Error())
			os.Exit(1)
		}
		secrets = append(secrets, secret)
	}

	if len(secrets) == 0 {
		fmt.Print("No systems of record found\n")
		return
	}

	for _, secret := range secrets {
		var value map[string]interface{}
		err := json.Unmarshal([]byte(*secret.Value), &value)
		if err != nil {
			log.Printf("failed to retrieve systems; %s", err.Error())
			os.Exit(1)
		}

		result, _ := json.MarshalIndent(value, "", "\t")
		fmt.Printf("%s\n", string(result))
	}
}

func init() {
	listBaselineSystemsCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	listBaselineSystemsCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")

	listBaselineSystemsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineSystemsCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of baseline workgroups to retrieve per page")
	listBaselineSystemsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
