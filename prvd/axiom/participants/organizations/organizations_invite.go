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

package organizations

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kthomas/go-pgputil"
	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	ident "github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/nchain"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var firstName string
var lastName string
var email string

var orgName string

var inviteBaselineWorkgroupOrganizationCmd = &cobra.Command{
	Use:   "invite",
	Short: "Invite an organization to a axiom workgroup",
	Long: `Invite an organization to participate in a axiom workgroup.
  
  A verifiable credential is issued which can then be distributed to the invited party out-of-band.`,
	Run: inviteOrganization,
}

func inviteOrganization(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInvite)
}

func inviteOrganizationRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}
	if firstName == "" {
		firstNamePrompt()
	}
	if lastName == "" {
		lastNamePrompt()
	}
	if email == "" {
		emailPrompt()
	}
	if orgName == "" {
		orgNamePrompt()
	}

	common.AuthorizeOrganizationContext(false)

	token, err := common.ResolveOrganizationToken()

	vaults, err := vault.ListVaults(*token.AccessToken, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to resolve vault for organization; %s", err.Error())
		os.Exit(1)
	}
	orgVaultID := vaults[0].ID.String()

	keys, err := vault.ListKeys(*token.AccessToken, orgVaultID, map[string]interface{}{
		"spec": "secp256k1",
	})
	if err != nil {
		log.Printf("failed to resolve secp256k1 key for organization; %s", err.Error())
		os.Exit(1)
	}
	secp256k1KeyAddress := keys[0].Address

	contracts, _ := nchain.ListContracts(*token.AccessToken, map[string]interface{}{
		"type": "organization-registry",
	})
	if err != nil {
		log.Printf("failed to resolve contract for organization; %s", err.Error())
		os.Exit(1)
	}
	orgRegistryAddress := contracts[0].Address

	if common.SubjectAccountID == "" {
		common.SubjectAccountID = common.SHA256(fmt.Sprintf("%s.%s", common.OrganizationID, common.WorkgroupID))
	}

	jwtParams := map[string]interface{}{
		"invitor_organization_address": secp256k1KeyAddress,
		"registry_contract_address":    orgRegistryAddress,
		"workgroup_id":                 common.WorkgroupID,
		"invitor_subject_account_id":   common.SubjectAccountID,
	}

	authorizedBearerToken := vendJWT(orgVaultID, jwtParams)

	wgID, _ := uuid.FromString(common.WorkgroupID)

	inviteParams := map[string]interface{}{
		"first_name":        firstName,
		"last_name":         lastName,
		"email":             email,
		"organization_name": orgName,
		"application_id":    common.WorkgroupID, // FIXME-- should be workgroup id
		"params": map[string]interface{}{
			"verifiable_credential":      authorizedBearerToken,
			"is_organization_invite":     true,
			"operator_separation_degree": common.Organization.Metadata.Workgroups[wgID].OperatorSeparationDegree + 1,
			"workgroup":                  common.Workgroup,
		},
	}

	if err := ident.CreateInvitation(*token.AccessToken, inviteParams); err != nil {
		log.Printf("failed to invite axiom workgroup user; %s", err.Error())
		os.Exit(1)
	}

	log.Printf("invited axiom workgroup organization: %s\n", orgName)
}

func vendJWT(vaultID string, params map[string]interface{}) string {
	keys, err := vault.ListKeys(common.OrganizationAccessToken, vaultID, map[string]interface{}{
		"spec": "RSA-4096",
	})
	if err != nil {
		log.Printf("failed to resolve RSA-4096 key for organization; %s", err.Error())
		os.Exit(1)
	}
	if len(keys) == 0 {
		log.Print("failed to resolve RSA-4096 key for organization")
		os.Exit(1)
	}
	key := keys[0]

	org, err := ident.GetOrganizationDetails(common.OrganizationAccessToken, common.OrganizationID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to vend JWT; %s", err.Error())
		os.Exit(1)
	}

	issuedAt := time.Now()

	claims := map[string]interface{}{
		"aud":   org.Metadata["messaging_endpoint"],
		"iat":   issuedAt.Unix(),
		"iss":   common.OrganizationID,
		"sub":   email,
		"axiom": params,
	}

	natsClaims, err := encodeJWTNatsClaims()
	if err != nil {
		log.Printf("failed to encode NATS claims in JWT; %s", err.Error())
		os.Exit(1)
	}
	if natsClaims != nil {
		claims["nats"] = natsClaims
	}

	publicKey, err := pgputil.DecodeRSAPublicKeyFromPEM([]byte(*key.PublicKey))
	if err != nil {
		log.Printf("failed to decode RSA public key from PEM; %s", err.Error())
		os.Exit(1)
	}

	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		log.Printf("failed to decode SSH public key for fingerprinting; %s", err.Error())
		os.Exit(1)
	}
	fingerprint := ssh.FingerprintLegacyMD5(sshPublicKey)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claims))
	jwtToken.Header["kid"] = fingerprint

	strToSign, err := jwtToken.SigningString()
	if err != nil {
		log.Printf("failed to generate JWT string for signing; %s", err.Error())
		os.Exit(1)
	}

	opts := map[string]interface{}{}
	if strings.HasPrefix(*key.Spec, "RSA-") {
		opts["algorithm"] = "RS256"
	}

	resp, err := vault.SignMessage(
		common.OrganizationAccessToken,
		key.VaultID.String(),
		key.ID.String(),
		hex.EncodeToString([]byte(strToSign)),
		opts,
	)
	if err != nil {
		log.Printf("WARNING: failed to sign JWT using vault key: %s; %s", key.ID, err.Error())
		os.Exit(1)
	}

	sigAsBytes, err := hex.DecodeString(*resp.Signature)
	if err != nil {
		log.Printf("failed to decode signature from hex; %s", err.Error())
		os.Exit(1)
	}

	encodedSignature := strings.TrimRight(base64.URLEncoding.EncodeToString(sigAsBytes), "=")
	return strings.Join([]string{strToSign, encodedSignature}, ".")
}

func encodeJWTNatsClaims() (map[string]interface{}, error) {
	publishAllow := make([]string, 0)
	publishDeny := make([]string, 0)

	subscribeAllow := make([]string, 0)
	subscribeDeny := make([]string, 0)

	var responsesMax *int
	var responsesTTL *time.Duration

	// subscribeAllow = append(subscribeAllow, "axiom.>")
	publishAllow = append(publishAllow, "axiom.>")

	var publishPermissions map[string]interface{}
	if len(publishAllow) > 0 || len(publishDeny) > 0 {
		publishPermissions = map[string]interface{}{}
		if len(publishAllow) > 0 {
			publishPermissions["allow"] = publishAllow
		}
		if len(publishDeny) > 0 {
			publishPermissions["deny"] = publishDeny
		}
	}

	var subscribePermissions map[string]interface{}
	if len(subscribeAllow) > 0 || len(subscribeDeny) > 0 {
		subscribePermissions = map[string]interface{}{}
		if len(subscribeAllow) > 0 {
			subscribePermissions["allow"] = subscribeAllow
		}
		if len(subscribeDeny) > 0 {
			subscribePermissions["deny"] = subscribeDeny
		}
	}

	var responsesPermissions map[string]interface{}
	if responsesMax != nil || responsesTTL != nil {
		responsesPermissions = map[string]interface{}{}
		if responsesMax != nil {
			responsesPermissions["max"] = responsesMax
		}
		if responsesTTL != nil {
			responsesPermissions["ttl"] = responsesTTL
		}
	}

	var permissions map[string]interface{}
	if publishPermissions != nil || subscribePermissions != nil || responsesPermissions != nil {
		permissions = map[string]interface{}{}
		if publishPermissions != nil {
			permissions["publish"] = publishPermissions
		}
		if subscribePermissions != nil {
			permissions["subscribe"] = subscribePermissions
		}
		if responsesPermissions != nil {
			permissions["responses"] = responsesPermissions
		}
	}

	var natsClaims map[string]interface{}
	if permissions != nil {
		natsClaims = map[string]interface{}{
			"permissions": permissions,
		}
	}

	return natsClaims, nil
}

func firstNamePrompt() {
	prompt := promptui.Prompt{
		Label: "Invitee First Name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("first name required")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	firstName = result
}

func lastNamePrompt() {
	prompt := promptui.Prompt{
		Label: "Invitee Last Name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("last name required")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	lastName = result
}

func emailPrompt() {
	prompt := promptui.Prompt{
		Label:    "Invitee Email",
		Validate: common.EmailValidation,
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	email = result
}

func orgNamePrompt() {
	prompt := promptui.Prompt{
		Label: "Invitee Organization Name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("organization name required")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	orgName = result
}

func init() {
	inviteBaselineWorkgroupOrganizationCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	inviteBaselineWorkgroupOrganizationCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	inviteBaselineWorkgroupOrganizationCmd.Flags().StringVar(&common.SubjectAccountID, "subject-account", "", "subject account identifier")
	inviteBaselineWorkgroupOrganizationCmd.Flags().StringVar(&firstName, "first-name", "", "first name of the invited participant")
	inviteBaselineWorkgroupOrganizationCmd.Flags().StringVar(&lastName, "last-name", "", "last name of the invited participant")
	inviteBaselineWorkgroupOrganizationCmd.Flags().StringVar(&email, "email", "", "email address of the invited participant")
	inviteBaselineWorkgroupOrganizationCmd.Flags().StringVar(&orgName, "organization-name", "", "name of the invited organization")

	// inviteBaselineWorkgroupOrganizationCmd.Flags().BoolVar(&managedTenant, "managed-tenant", false, "if set, the invited participant is authorized to leverage operator-provided infrastructure")
	// inviteBaselineWorkgroupOrganizationCmd.Flags().IntVar(&permissions, "permissions", 0, "permissions for invited participant")
	inviteBaselineWorkgroupOrganizationCmd.Flags().BoolVarP(&Optional, "Optional", "", false, "List all the Optional flags")
}
