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
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-cli/prvd/organizations"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inviteJWT string

var mode string

var firstName string
var lastName string

var email string
var password string

var orgName string
var orgDescription string

var joinBaselineWorkgroupCmd = &cobra.Command{
	Use:   "join",
	Short: "Join a baseline workgroup",
	Long:  `Join a baseline workgroup by accepting the invite.`,
	Run:   joinWorkgroup,
}

func joinWorkgroup(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepJoin)
}

func joinWorkgroupRun(cmd *cobra.Command, args []string) {
	if inviteJWT == "" {
		jwtPrompt()
	}

	decodedTokenData, err := parseJWT(inviteJWT)
	if err != nil {
		fmt.Printf("failed to accept invitation; %s", err.Error())
	}

	if mode != "login" && mode != "signup" {
		modePrompt()
	}

	if mode == "login" {
		loginPrompt(*decodedTokenData.Email)
	} else if mode == "signup" {
		signupPrompt(*decodedTokenData.FirstName, *decodedTokenData.LastName, *decodedTokenData.Email)
	}

	if mode == "login" {
		resp, err := ident.Authenticate(email, password, inviteJWT)
		if err != nil {
			fmt.Printf("failed to accept invite; %s", err.Error())
			os.Exit(1)
		}

		if resp.Token.AccessToken != nil && resp.Token.RefreshToken != nil {
			common.CacheAccessRefreshToken(resp.Token, nil)
		} else if resp.Token.Token != nil {
			viper.Set(common.AccessTokenConfigKey, *resp.Token.Token)
			viper.WriteConfig()
		}
	} else if mode == "signup" {
		createUserParams := map[string]interface{}{
			"first_name": firstName,
			"last_name":  lastName,
			"email":      email,
			"password":   password,
		}

		if decodedTokenData.Params.IsOrganizationUserInvite {
			createUserParams["invitation_token"] = inviteJWT
		}

		if _, err := ident.CreateUser("", createUserParams); err != nil {
			fmt.Printf("failed to accept invite; %s", err.Error())
			os.Exit(1)
		}

		resp, err := ident.Authenticate(email, password, "")
		if err != nil {
			fmt.Printf("failed to accept invite; %s", err.Error())
			os.Exit(1)
		}

		if resp.Token.AccessToken != nil && resp.Token.RefreshToken != nil {
			common.CacheAccessRefreshToken(resp.Token, nil)
		} else if resp.Token.Token != nil {
			viper.Set(common.AccessTokenConfigKey, *resp.Token.Token)
			viper.WriteConfig()
		}
	}

	if decodedTokenData.Params.IsOrganizationInvite {
		orgPrompt(*decodedTokenData.OrganizationName)

		if !hasAgreedToTermsOfService {
			if ok := common.RequireTermsOfServiceAgreement(); !ok {
				fmt.Print("failed to accept invite; must accept the terms of agreement")
				os.Exit(1)
			}
		}
		if !hasAgreedToPrivacyPolicy {
			if ok := common.RequirePrivacyPolicyAgreement(); !ok {
				fmt.Print("failed to accept invite; must accept the privacy policy")
				os.Exit(1)
			}
		}

		orgParams := map[string]interface{}{
			"name": orgName,
			"metadata": map[string]interface{}{
				"workgroups": map[uuid.UUID]interface{}{
					*decodedTokenData.ApplicationID: map[string]interface{}{
						"operator_separation_degree": decodedTokenData.Params.OperatorSeparationDegree,
						"system_secret_ids":          make([]string, 0),
						"vault_id":                   nil,
					},
				},
			},
			"invitation_token": inviteJWT,
		}

		if orgDescription != "" {
			orgParams["description"] = orgDescription
		}

		org, err := ident.CreateOrganization(common.RequireUserAccessToken(), orgParams)
		if err != nil {
			fmt.Printf("failed to accept invite; %s", err.Error())
			os.Exit(1)
		}

		common.OrganizationID = *org.ID

		common.AuthorizeOrganizationContext(true)

		token := common.RequireOrganizationToken()

		common.RequireOrganizationVault()

		vaults, err := vault.ListVaults(token, map[string]interface{}{})
		if err != nil {
			log.Printf("failed to accept invite; %s", err.Error())
			os.Exit(1)
		}

		orgVault := vaults[0]
		if orgVault == nil {
			log.Print("failed to accept invite; failed to fetch organization vault; no vaults found")
			os.Exit(1)
		}

		requireOrganizationKeys()

		secp256k1Key, err := vault.FetchKey(token, common.VaultID, secp256k1KeyID)
		if err != nil {
			fmt.Printf("failed to initialize baseline workgroup: %s", err.Error())
			os.Exit(1)
		}

		termsNow := time.Now()
		termsString := termsNow.Format(time.RFC3339)
		termsTimestampHex := hex.EncodeToString([]byte(termsString))
		termsTimestampHash := common.SHA256(termsTimestampHex)

		termsSig, err := vault.SignMessage(token, common.VaultID, secp256k1KeyID, termsTimestampHash, map[string]interface{}{})
		if err != nil {
			fmt.Printf("failed to initialize baseline workgroup: %s", err.Error())
			os.Exit(1)
		}

		privacyNow := time.Now()
		privacyString := privacyNow.Format(time.RFC3339)
		privacyTimestampHex := hex.EncodeToString([]byte(privacyString))
		privacyTimestampHash := common.SHA256(privacyTimestampHex)

		privacySig, err := vault.SignMessage(token, common.VaultID, secp256k1KeyID, privacyTimestampHash, map[string]interface{}{})
		if err != nil {
			fmt.Printf("failed to initialize baseline workgroup: %s", err.Error())
			os.Exit(1)
		}

		// TODO-- configure system for organization participant

		var localOrg organizations.Organization
		raw, _ := json.Marshal(org)
		json.Unmarshal(raw, &localOrg)

		localOrg.Metadata.Address = *secp256k1Key.Address
		localOrg.Metadata.Workgroups[*decodedTokenData.ApplicationID] = &organizations.OrganizationWorkgroupMetadata{
			OperatorSeparationDegree: decodedTokenData.Params.OperatorSeparationDegree,
			VaultID:                  &orgVault.ID,
			SystemSecretIDs:          make([]*uuid.UUID, 0),
			TOS: &organizations.WorkgroupMetadataLegal{
				AgreedAt:  &termsNow,
				Signature: termsSig.Signature,
			},
			Privacy: &organizations.WorkgroupMetadataLegal{
				AgreedAt:  &privacyNow,
				Signature: privacySig.Signature,
			},
		}

		var orgInterface map[string]interface{}
		raw, _ = json.Marshal(localOrg)
		json.Unmarshal(raw, &orgInterface)

		err = ident.UpdateOrganization(token, common.OrganizationID, orgInterface)
		if err != nil {
			log.Printf("failed to accept invite; %s", err.Error())
			os.Exit(1)
		}

		refresh := common.RequireOrganizationRefreshToken()

		subjectAccountParams := map[string]interface{}{
			"metadata": map[string]interface{}{
				"organization_id":            common.OrganizationID,
				"organization_address":       *secp256k1Key.Address,
				"organization_refresh_token": refresh,
				"workgroup_id":               *decodedTokenData.ApplicationID,
				"registry_contract_address":  *secp256k1Key.Address,
				"network_id":                 *decodedTokenData.Params.Workgroup.NetworkID,
			},
		}

		_, err = baseline.CreateWorkgroup(token, map[string]interface{}{
			"subject_account_params": subjectAccountParams,
			"token":                  *decodedTokenData.Params.AuthorizedBearerToken,
		})
		if err != nil {
			log.Printf("failed to accept invite; %s", err.Error())
			os.Exit(1)
		}
	}

	fmt.Print("successfully accepted invitation\n") // TODO-- elaborate
}

func jwtPrompt() {
	prompt := promptui.Prompt{
		Label: "Verifiable Credential (Invite JWT)",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	inviteJWT = result
}

func parseJWT(token string) (*baseline.DataClaims, error) {
	claims := &baseline.InviteClaims{}

	var jwtParser jwt.Parser
	jwtToken, _, err := jwtParser.ParseUnverified(token, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT; %s", err.Error())
	}

	var data baseline.InviteClaims
	raw, err := json.Marshal(jwtToken.Claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT; %s", err.Error())
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT; %s", err.Error())
	}

	return data.PRVD.Data, nil
}

func modePrompt() {
	acceptInviteModes := make([]string, 2)
	acceptInviteModes[0] = "signup"
	acceptInviteModes[1] = "login"

	selectPrompt := promptui.Select{
		Label: "Signup or login",
		Items: acceptInviteModes,
	}

	i, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("failed to accept invitation; %s", err.Error())
		os.Exit(1)
	}

	mode = acceptInviteModes[i]
}

func loginPrompt(defaultEmail string) {
	if email == "" {
		prompt := promptui.Prompt{
			Label:    "Email",
			Validate: common.EmailValidation,
			Default:  defaultEmail,
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		email = result
	}

	if password == "" {
		password = common.FreeInput("Password", "", common.MandatoryValidation)
	}
}

func signupPrompt(defaultFirst, defaultLast, defaultEmail string) {
	if firstName == "" {
		prompt := promptui.Prompt{
			Label:    "First Name",
			Validate: common.MandatoryValidation,
			Default:  defaultFirst,
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		firstName = result
	}

	if lastName == "" {
		prompt := promptui.Prompt{
			Label:    "Last Name",
			Validate: common.MandatoryValidation,
			Default:  defaultLast,
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		lastName = result
	}

	if email == "" {
		prompt := promptui.Prompt{
			Label:    "Email",
			Validate: common.EmailValidation,
			Default:  defaultEmail,
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		email = result
	}

	if password == "" {
		password = common.FreeInput("Password", "", common.MandatoryValidation)
	}
}

func orgPrompt(defaultName string) {
	if orgName == "" {
		prompt := promptui.Prompt{
			Label:    "Organization Name",
			Validate: common.MandatoryValidation,
			Default:  defaultName,
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		orgName = result
	}

	if orgDescription == "" {
		prompt := promptui.Prompt{
			Label: "Organization Description",
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		description = result
	}
}

func init() {
	joinBaselineWorkgroupCmd.Flags().StringVar(&inviteJWT, "jwt", "", "JWT invitation token received from the inviting counterparty")

	joinBaselineWorkgroupCmd.Flags().StringVar(&mode, "mode", "", "signup or login to accept invitation")

	joinBaselineWorkgroupCmd.Flags().StringVar(&firstName, "first-name", "", "first name of created user to accept invitation")
	joinBaselineWorkgroupCmd.Flags().StringVar(&lastName, "last-name", "", "last name of created user to accept invitation")

	joinBaselineWorkgroupCmd.Flags().StringVar(&email, "email", "", "email of created user to accept invitation")
	joinBaselineWorkgroupCmd.Flags().StringVar(&password, "password", "", "password of created user to accept invitation")

	joinBaselineWorkgroupCmd.Flags().StringVar(&orgName, "organization-name", "", "organization name of invited organization")
	joinBaselineWorkgroupCmd.Flags().StringVar(&orgDescription, "organization-description", "", "organization description of invited organization")

	joinBaselineWorkgroupCmd.Flags().BoolVarP(&hasAgreedToTermsOfService, "terms", "", false, "accept the terms of service (https://provide.services/terms)")
	joinBaselineWorkgroupCmd.Flags().BoolVarP(&hasAgreedToPrivacyPolicy, "privacy", "", false, "accept the privacy policy (https://provide.services/privacy-policy)")

	joinBaselineWorkgroupCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")

}
