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

	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/axiom"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/nchain"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/spf13/cobra"
)

var orgDomain string

var Optional bool
var paginate bool

var initBaselineSubjectAccountCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize axiom subject account",
	Long:  `Initialize and configure a new axiom subject account`,
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

	if common.Organization.Metadata != nil && common.Organization.Metadata.Domain != "" {
		orgDomain = common.Organization.Metadata.Domain
	} else if orgDomain == "" {
		orgDomainPrompt()
	}

	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()

	contracts, err := nchain.ListContracts(*token.AccessToken, map[string]interface{}{
		"type": "organization-registry",
	})
	if err != nil {
		fmt.Printf("failed to create subject account; %s", err.Error())
		os.Exit(1)
	}

	if len(contracts) == 0 {
		fmt.Println("failed to create subject account; failed to resolve organization registry contract")
		os.Exit(1)
	}

	if len(contracts) > 1 {
		fmt.Println("failed to create subject account; resolved ambiguous organization registry contracts")
		os.Exit(1)
	}

	sa, err := axiom.CreateSubjectAccount(*token.AccessToken, common.OrganizationID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"organization_id":            common.OrganizationID,
			"organization_address":       common.Organization.Metadata.Address,
			"organization_refresh_token": *token.RefreshToken,
			"workgroup_id":               common.WorkgroupID,
			"registry_contract_address":  *contracts[0].Address,
			"network_id":                 common.NetworkID,
			"organization_domain":        orgDomain,
		},
	})
	if err != nil {
		log.Printf("failed to initialize axiom subject account; %s", err.Error())
		os.Exit(1)
	}

	// TODO-- make utility function to DRY this up
	vaultID := common.Organization.Metadata.Workgroups[common.Workgroup.ID].VaultID
	systemIDs := common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs

	isOperator := common.Organization.Metadata.Workgroups[common.Workgroup.ID].OperatorSeparationDegree == 0
	if isOperator {
		vaultID = common.Workgroup.Config.VaultID
		systemIDs = common.Workgroup.Config.SystemSecretIDs
	}

	if len(systemIDs) > 0 && vaultID != nil {
		for _, secretID := range systemIDs {
			secret, err := vault.FetchSecret(*token.AccessToken, vaultID.String(), secretID.String(), map[string]interface{}{})
			if err != nil {
				log.Printf("failed to initialize axiom subject account; %s", err.Error())
				os.Exit(1)
			}

			var systemParams map[string]interface{}
			err = json.Unmarshal([]byte(*secret.Value), &systemParams)
			if err != nil {
				log.Printf("failed to initialize axiom subject account; %s", err.Error())
				os.Exit(1)
			}

			if _, err := axiom.CreateSystem(*token.AccessToken, common.WorkgroupID, systemParams); err != nil {
				log.Printf("failed to initialize axiom subject account; %s", err.Error())
				os.Exit(1)
			}

			if err := vault.DeleteSecret(*token.AccessToken, vaultID.String(), secretID.String()); err != nil {
				log.Printf("failed to initialize axiom subject account; %s", err.Error())
				os.Exit(1)
			}
		}

		common.Organization.Metadata.Domain = orgDomain
		common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs = make([]*uuid.UUID, 0)

		var organizationParams map[string]interface{}
		raw, _ := json.Marshal(common.Organization)
		json.Unmarshal(raw, &organizationParams)

		if err := ident.UpdateOrganization(*token.AccessToken, common.OrganizationID, organizationParams); err != nil {
			log.Printf("failed to initialize axiom subject account; %s", err.Error())
			os.Exit(1)
		}

		if isOperator {
			common.Workgroup.Config.SystemSecretIDs = make([]*uuid.UUID, 0)

			var workgroupParams map[string]interface{}
			raw, _ := json.Marshal(common.Workgroup)
			json.Unmarshal(raw, &workgroupParams)

			if err := axiom.UpdateWorkgroup(*token.AccessToken, common.WorkgroupID, workgroupParams); err != nil {
				log.Printf("failed to initialize axiom subject account; %s", err.Error())
				os.Exit(1)
			}
		}

		fmt.Printf("successfully saved systems of records to subject account %s\n", *sa.ID)
	}

	result, _ := json.MarshalIndent(sa, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func orgDomainPrompt() {
	prompt := promptui.Prompt{
		Label: "Organization Domain",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("org domain cannot be empty")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	orgDomain = result
}

func init() {
	initBaselineSubjectAccountCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineSubjectAccountCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	initBaselineSubjectAccountCmd.Flags().StringVar(&common.NetworkID, "network", "", "nchain network id of the axiom mainnet to use for this workgroup")

	initBaselineSubjectAccountCmd.Flags().StringVar(&orgDomain, "organization-domain", "", "organization domain to use for this subject account, if it is not set on the organization")

	initBaselineSubjectAccountCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
