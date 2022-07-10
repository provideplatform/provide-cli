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
	"log"
	"os"
	"strings"

	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-cli/prvd/organizations"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/nchain"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/spf13/cobra"
)

const defaultNChainBaselineNetworkID = "66d44f30-9092-4182-a3c4-bc02736d6ae5"

var name string

var vaultID string
var babyJubJubKeyID string
var secp256k1KeyID string
var hdwalletID string
var rsa4096KeyID string
var Optional bool
var paginate bool

var initBaselineWorkgroupCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize baseline workgroup",
	Long:  `Initialize and configure a new baseline workgroup`,
	Run:   initWorkgroup,
}

func AuthorizeApplicationContext() {
	// common.AuthorizeApplicationContext()
	_, err := nchain.CreateWallet(common.ApplicationAccessToken, map[string]interface{}{
		"purpose": 44,
	})
	if err != nil {
		log.Printf("failed to initialize HD wallet; %s", err.Error())
		os.Exit(1)
	}
}
func initWorkgroup(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func initWorkgroupRun(cmd *cobra.Command, args []string) {
	if name == "" {
		namePrompt()
	}
	if common.NetworkID == "" {
		common.RequirePublicNetwork()
	}
	common.AuthorizeOrganizationContext(true)

	token := common.RequireOrganizationToken()

	vaults, err := vault.ListVaults(token, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	orgVault := vaults[0]
	if orgVault == nil {
		log.Print("failed to initialize baseline workgroup; failed to fetch organization vault; no vaults found")
		os.Exit(1)
	}

	wg, err := baseline.CreateWorkgroup(token, map[string]interface{}{
		"name":       name,
		"network_id": common.NetworkID,
		"config": map[string]interface{}{
			"vault_id": orgVault.ID.String(),
		},
		"type": "baseline",
	})
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	org, err := ident.GetOrganizationDetails(token, common.OrganizationID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	var organization organizations.Organization
	raw, _ := json.Marshal(org)
	json.Unmarshal(raw, &organization)

	organization.Metadata = &organizations.OrganizationMetadata{
		Workgroups: map[uuid.UUID]*organizations.OrganizationWorkgroupMetadata{},
	}

	organization.Metadata.Workgroups[wg.ID] = &organizations.OrganizationWorkgroupMetadata{
		OperatorSeparationDegree: uint32(0),
		VaultID:                  &orgVault.ID,
	}

	var orgInterface map[string]interface{}
	raw, _ = json.Marshal(organization)
	json.Unmarshal(raw, &orgInterface)

	err = ident.UpdateOrganization(token, common.OrganizationID, orgInterface)
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	common.WorkgroupID = wg.ID.String()

	common.InitWorkgroupContract()

	common.RequireOrganizationVault()
	requireOrganizationKeys()
	// common.RegisterWorkgroupOrganization(wg.ID.String())

	refresh := common.RequireOrganizationRefreshToken()

	secp256k1Key, err := vault.FetchKey(token, common.VaultID, secp256k1KeyID)
	if err != nil {
		fmt.Printf("failed to initialize baseline workgroup: %s", err.Error())
		os.Exit(1)
	}

	sa, err := baseline.CreateSubjectAccount(token, common.OrganizationID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"organization_id":            common.OrganizationID,
			"organization_address":       *secp256k1Key.Address,
			"organization_refresh_token": refresh,
			"workgroup_id":               common.WorkgroupID,
			"registry_contract_address":  *secp256k1Key.Address,
			"network_id":                 common.NetworkID,
		},
	})

	//common.RequireOrganizationEndpoints(nil)

	log.Printf("initialized baseline workgroup: %s; subject account id: %s", wg.ID, *sa.ID)
}

func requireOrganizationKeys() {
	var key *vault.Key
	var err error

	key, err = common.RequireOrganizationKeypair("babyJubJub")
	if err == nil {
		babyJubJubKeyID = key.ID.String()
	}

	key, err = common.RequireOrganizationKeypair("secp256k1")
	if err == nil {
		secp256k1KeyID = key.ID.String()
	}

	key, err = common.RequireOrganizationKeypair("BIP39")
	if err == nil {
		hdwalletID = key.ID.String()
	}

	key, err = common.RequireOrganizationKeypair("RSA-4096")
	if err == nil {
		rsa4096KeyID = key.ID.String()
	}
}

func namePrompt() {
	prompt := promptui.Prompt{
		Label: "Workgroup Name",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	name = result
}

func organizationAuthPrompt(target string) {
	prompt := promptui.Prompt{
		IsConfirm: true,
		Label:     fmt.Sprintf("Authorize access/refresh token for %s?", target),
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	if strings.ToLower(result) == "y" {
		common.AuthorizeOrganizationContext(true)
	}
}

func init() {
	initBaselineWorkgroupCmd.Flags().StringVar(&name, "name", "", "name of the baseline workgroup")
	initBaselineWorkgroupCmd.Flags().StringVar(&common.NetworkID, "network", "", "nchain network id of the baseline mainnet to use for this workgroup")
	initBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineWorkgroupCmd.Flags().StringVar(&common.MessagingEndpoint, "endpoint", "", "public messaging endpoint used for sending and receiving protocol messages")
	initBaselineWorkgroupCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
