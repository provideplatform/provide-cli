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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/nchain"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/spf13/cobra"
)

const defaultNChainBaselineNetworkID = "66d44f30-9092-4182-a3c4-bc02736d6ae5"

var name string
var description string

var hasAgreedToTermsOfService bool
var hasAgreedToPrivacyPolicy bool

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
	if _, err := nchain.CreateWallet(common.ApplicationAccessToken, map[string]interface{}{
		"purpose": 44,
	}); err != nil {
		log.Printf("failed to initialize HD wallet; %s", err.Error())
		os.Exit(1)
	}
}
func initWorkgroup(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func initWorkgroupRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if name == "" {
		namePrompt()
	}
	if description == "" {
		descriptionPrompt()
	}
	if common.NetworkID == "" {
		common.RequireL1Network()
	}
	if common.L2NetworkID == "" {
		common.RequireL2Network()
	}
	if !hasAgreedToTermsOfService {
		if ok := common.RequireTermsOfServiceAgreement(); !ok {
			fmt.Print("failed to initialize baseline workgroup; must accept the terms of agreement")
			os.Exit(1)
		}
	}
	if !hasAgreedToPrivacyPolicy {
		if ok := common.RequirePrivacyPolicyAgreement(); !ok {
			fmt.Print("failed to initialize baseline workgroup; must accept the privacy policy")
			os.Exit(1)
		}
	}

	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	vaults, err := vault.ListVaults(*token.AccessToken, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	orgVault := vaults[0]
	if orgVault == nil {
		log.Print("failed to initialize baseline workgroup; failed to fetch organization vault; no vaults found")
		os.Exit(1)
	}

	params := map[string]interface{}{
		"name":       name,
		"network_id": common.NetworkID,
		"config": map[string]interface{}{
			"vault_id":          orgVault.ID.String(),
			"l2_network_id":     common.L2NetworkID,
			"system_secret_ids": make([]*uuid.UUID, 0),
		},
		"type": "baseline",
	}

	if description != "" {
		params["description"] = description
	}

	wg, err := baseline.CreateWorkgroup(*token.AccessToken, params)
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	if err := common.RequireOrganizationVault(); err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	if err := requireOrganizationKeys(); err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	secp256k1Key, err := vault.FetchKey(*token.AccessToken, common.VaultID, secp256k1KeyID)
	if err != nil {
		fmt.Printf("failed to initialize baseline workgroup: %s", err.Error())
		os.Exit(1)
	}

	termsNow := time.Now()
	termsString := termsNow.Format(time.RFC3339)
	termsTimestampHex := hex.EncodeToString([]byte(termsString))
	termsTimestampHash := common.SHA256(termsTimestampHex)

	termsSig, err := vault.SignMessage(*token.AccessToken, common.VaultID, secp256k1KeyID, termsTimestampHash, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to initialize baseline workgroup: %s", err.Error())
		os.Exit(1)
	}

	privacyNow := time.Now()
	privacyString := privacyNow.Format(time.RFC3339)
	privacyTimestampHex := hex.EncodeToString([]byte(privacyString))
	privacyTimestampHash := common.SHA256(privacyTimestampHex)

	privacySig, err := vault.SignMessage(*token.AccessToken, common.VaultID, secp256k1KeyID, privacyTimestampHash, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to initialize baseline workgroup: %s", err.Error())
		os.Exit(1)
	}

	if common.Organization.Metadata == nil {
		common.Organization.Metadata = &common.OrganizationMetadata{
			Address:    *secp256k1Key.Address,
			Workgroups: map[uuid.UUID]*common.OrganizationWorkgroupMetadata{},
		}
	} else if common.Organization.Metadata.Workgroups == nil {
		common.Organization.Metadata.Workgroups = map[uuid.UUID]*common.OrganizationWorkgroupMetadata{}
	}

	common.Organization.Metadata.Workgroups[wg.ID] = &common.OrganizationWorkgroupMetadata{
		OperatorSeparationDegree: uint32(0),
		VaultID:                  &orgVault.ID,
		SystemSecretIDs:          make([]*uuid.UUID, 0),
		TOS: &common.WorkgroupMetadataLegal{
			AgreedAt:  &termsNow,
			Signature: termsSig.Signature,
		},
		Privacy: &common.WorkgroupMetadataLegal{
			AgreedAt:  &privacyNow,
			Signature: privacySig.Signature,
		},
	}

	var orgInterface map[string]interface{}
	raw, _ := json.Marshal(common.Organization)
	json.Unmarshal(raw, &orgInterface)

	if err := ident.UpdateOrganization(*token.AccessToken, common.OrganizationID, orgInterface); err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	common.WorkgroupID = wg.ID.String()

	common.InitWorkgroupContract()

	sa, err := baseline.CreateSubjectAccount(*token.AccessToken, common.OrganizationID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"organization_id":            common.OrganizationID,
			"organization_address":       *secp256k1Key.Address,
			"organization_refresh_token": *token.AccessToken,
			"workgroup_id":               common.WorkgroupID,
			"registry_contract_address":  *secp256k1Key.Address,
			"network_id":                 common.NetworkID,
		},
	})

	//common.RequireOrganizationEndpoints(nil)
	result, _ := json.MarshalIndent(wg, "", "\t")
	fmt.Printf("%s\nsubject account id: %s\n", string(result), *sa.ID)
}

func requireOrganizationKeys() error {
	var key *vault.Key
	var err error

	key, err = common.RequireOrganizationKeypair("babyJubJub")
	if err != nil {
		return err
	}
	babyJubJubKeyID = key.ID.String()

	key, err = common.RequireOrganizationKeypair("secp256k1")
	if err != nil {
		return err
	}
	secp256k1KeyID = key.ID.String()

	key, err = common.RequireOrganizationKeypair("BIP39")
	if err != nil {
		return err
	}
	hdwalletID = key.ID.String()

	key, err = common.RequireOrganizationKeypair("RSA-4096")
	if err != nil {
		return err
	}
	rsa4096KeyID = key.ID.String()

	return nil
}

func namePrompt() {
	prompt := promptui.Prompt{
		Label: "Workgroup Name",
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
		return
	}

	name = result
}

func descriptionPrompt() {
	prompt := promptui.Prompt{
		Label: "Workgroup Description",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	description = result
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
	initBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineWorkgroupCmd.Flags().StringVar(&name, "name", "", "name of the baseline workgroup")
	initBaselineWorkgroupCmd.Flags().StringVar(&description, "description", "", "description of the baseline workgroup")
	initBaselineWorkgroupCmd.Flags().StringVar(&common.NetworkID, "network", "", "nchain network id of the baseline mainnet to use for this workgroup")
	initBaselineWorkgroupCmd.Flags().StringVar(&common.L2NetworkID, "l2", "", "nchain l2 network id of the baseline layer 2 to use for this workgroup")
	initBaselineWorkgroupCmd.Flags().BoolVarP(&hasAgreedToTermsOfService, "terms", "", false, "accept the terms of service (https://provide.services/terms)")
	initBaselineWorkgroupCmd.Flags().BoolVarP(&hasAgreedToPrivacyPolicy, "privacy", "", false, "accept the privacy policy (https://provide.services/privacy-policy)")

	initBaselineWorkgroupCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
