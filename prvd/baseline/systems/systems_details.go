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

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var vaultID string
var systemID string

var detailBaselineSystemCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific baseline system",
	Long:  `Retrieve details for a specific baseline system of record by identifier, scoped to the authorized API token`,
	Run:   fetchSystemDetails,
}

func fetchSystemDetails(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepDetails)
}

func fetchSystemDetailsRun(cmd *cobra.Command, args []string) {
	if err := common.RequireOrganization(); err != nil {
		fmt.Printf("failed to retrive system details; %s", err.Error())
		os.Exit(1)
	}

	if err := common.RequireWorkgroup(); err != nil {
		fmt.Printf("failed to retrive system details; %s", err.Error())
		os.Exit(1)
	}

	localVaultID := common.Organization.Metadata.Workgroups[common.Workgroup.ID].VaultID
	localSystemIDs := common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs

	isOperator := common.Organization.Metadata.Workgroups[common.Workgroup.ID].OperatorSeparationDegree == 0
	if isOperator {
		localVaultID = common.Workgroup.Config.VaultID
		localSystemIDs = common.Workgroup.Config.SystemSecretIDs
	}

	if vaultID != "" && vaultID != localVaultID.String() {
		fmt.Print("failed to retrieve system details; invalid vault id")
		os.Exit(1)
	}

	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to retrieve systems; %s", err.Error())
		os.Exit(1)
	}

	var system vault.Secret

	secrets := make([]*vault.Secret, 0)
	secretOpts := make([]string, 0)

	for _, secretID := range localSystemIDs {
		secret, err := vault.FetchSecret(*token.AccessToken, localVaultID.String(), secretID.String(), map[string]interface{}{})
		if err != nil {
			log.Printf("failed to retrieve systems; %s", err.Error())
			os.Exit(1)
		}

		secrets = append(secrets, secret)
		secretOpts = append(secretOpts, *secret.Description)
	}

	if systemID != "" {
		for _, s := range secrets {
			if systemID == s.ID.String() {
				system = *s
			}
		}

		if system.VaultID == nil {
			fmt.Print("failed to retrieve system details; invalid system id")
			os.Exit(1)
		}
	} else {
		prompt := promptui.Select{
			Label: "Select System",
			Items: secretOpts,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to retrieve system details; %s", err.Error())
			os.Exit(1)
		}

		system = *secrets[i]
	}

	var value map[string]interface{}
	err = json.Unmarshal([]byte(*system.Value), &value)
	if err != nil {
		log.Printf("failed to retrieve system details; %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(value, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func init() {
	detailBaselineSystemCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	detailBaselineSystemCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	detailBaselineSystemCmd.Flags().StringVar(&vaultID, "vault", "", "vault identifier")
	detailBaselineSystemCmd.Flags().StringVar(&systemID, "system", "", "system identifier")
}
