package participants

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
	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	ident "github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/api/nchain"
	"github.com/provideservices/provide-go/api/vault"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var name string
var email string
var permissions int
var invitorAddress string
var registryContractAddress string
var managedTenant bool
var Optional bool

var inviteBaselineWorkgroupParticipantCmd = &cobra.Command{
	Use:   "invite",
	Short: "Invite an organization to a baseline workgroup",
	Long: `Invite an organization to participate in a baseline workgroup.

A verifiable credential is issued which can then be distributed to the invited party out-of-band.`,
	Run: inviteParticipant,
}

func inviteParticipant(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInvite)
}

func inviteParticipantRun(cmd *cobra.Command, args []string) {
	if name == "" {
		namePrompt()
	}
	if email == "" {
		emailPrompt()
	}
	common.AuthorizeApplicationContext()
	common.AuthorizeOrganizationContext(false)

	vaults, err := vault.ListVaults(common.OrganizationAccessToken, map[string]interface{}{
		"organization_id": common.OrganizationID,
	})
	if err != nil {
		log.Printf("failed to resolve vault for organization; %s", err.Error())
		os.Exit(1)
	}

	keys, err := vault.ListKeys(common.OrganizationAccessToken, vaults[0].ID.String(), map[string]interface{}{
		"spec": "secp256k1",
	})
	if err != nil {
		log.Printf("failed to resolve secp256k1 key for organization; %s", err.Error())
		os.Exit(1)
	}

	contracts, _ := nchain.ListContracts(common.ApplicationAccessToken, map[string]interface{}{
		"type": "organization-registry",
	})
	if err != nil {
		log.Printf("failed to resolve contract for organization; %s", err.Error())
		os.Exit(1)
	}
	// TODO ... if no address prompt the user to get one
	invitorAddress := keys[0].Address
	registryContractAddress := contracts[0].Address

	params := map[string]interface{}{
		"invitor_organization_address": invitorAddress,
		"registry_contract_address":    registryContractAddress,
		"workgroup_id":                 common.ApplicationID,
	}

	authorizedBearerToken := vendJWT(vaults[0].ID.String(), params)
	params["authorized_bearer_token"] = authorizedBearerToken

	if common.OrganizationID != "" {
		params["organization_id"] = common.OrganizationID
	}

	if name != "" {
		params["organization_name"] = name
	}

	if permissions != 0 {
		params["permissions"] = permissions
	}

	// FIXME-- authorize the organization to act on behalf of this application when sending an invite
	err = ident.CreateInvitation(common.ApplicationAccessToken, map[string]interface{}{
		"application_id": common.ApplicationID,
		"email":          email,
		"params":         params,
	})
	if err != nil {
		log.Printf("failed to invite baseline workgroup participants; %s", err.Error())
		os.Exit(1)
	}

	log.Printf("invited baseline workgroup participant: %s\n\n\t%s", email, authorizedBearerToken)
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
		"aud":      org.Metadata["messaging_endpoint"],
		"iat":      issuedAt.Unix(),
		"iss":      fmt.Sprintf("organization:%s", common.OrganizationID),
		"sub":      email,
		"baseline": params,
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

	// subscribeAllow = append(subscribeAllow, "baseline.>")
	publishAllow = append(publishAllow, "baseline.>")

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

func namePrompt() {
	prompt := promptui.Prompt{
		Label: "Invitee Name",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	name = result
}

func emailPrompt() {
	prompt := promptui.Prompt{
		Label: "Invitee Email",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	email = result
}

func init() {
	inviteBaselineWorkgroupParticipantCmd.Flags().StringVar(&common.ApplicationID, "workgroup", "", "workgroup identifier")
	inviteBaselineWorkgroupParticipantCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	inviteBaselineWorkgroupParticipantCmd.Flags().StringVar(&name, "name", "", "name of the invited participant")
	inviteBaselineWorkgroupParticipantCmd.Flags().StringVar(&email, "email", "", "email address of the invited participant")
	inviteBaselineWorkgroupParticipantCmd.Flags().BoolVar(&managedTenant, "managed-tenant", false, "if set, the invited participant is authorized to leverage operator-provided infrastructure")
	inviteBaselineWorkgroupParticipantCmd.Flags().IntVar(&permissions, "permissions", 0, "permissions for invited participant")
	inviteBaselineWorkgroupParticipantCmd.Flags().BoolVarP(&Optional, "Optional", "", false, "List all the Optional flags")
}
