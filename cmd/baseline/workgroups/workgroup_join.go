package workgroups

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/kthomas/go-pgputil"
	"github.com/provideservices/provide-cli/cmd/common"
	ident "github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/common/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// InviteClaims represent JWT invitation claims
type InviteClaims struct {
	jwt.MapClaims
	Baseline *BaselineClaims `json:"baseline"`
}

// BaselineClaims represent JWT claims encoded within the invite token
type BaselineClaims struct {
	InvitorOrganizationAddress *string `json:"invitor_organization_address"`
	RegistryContractAddress    *string `json:"registry_contract_address"`
	WorkgroupID                *string `json:"workgroup_id"`
}

var inviteJWT string
var jwtKeypairs map[string]*util.JWTKeypair

var joinBaselineWorkgroupCmd = &cobra.Command{
	Use:   "join",
	Short: "Join a baseline workgroup",
	Long:  `Join a baseline workgroup by accepting the invite.`,
	Run:   joinWorkgroup,
}

func joinWorkgroup(cmd *cobra.Command, args []string) {
	requirePublicJWTVerifiers() // FIXME...
	claims := parseJWT(inviteJWT)
	log.Printf("resolved baseline claims containing invitation for workgroup: %s", *claims.Baseline.WorkgroupID)

	common.ApplicationID = *claims.Baseline.WorkgroupID
	authorizeOrganizationContext()
	authorizeApplicationContext()

	// initWorkgroupContract()

	requireOrganizationVault()
	requireOrganizationKeys()
	requireOrganizationMessagingEndpoint()
	registerWorkgroupOrganization(common.ApplicationID)
}

func parseJWT(token string) *InviteClaims {
	claims := &InviteClaims{}

	var jwtParser jwt.Parser
	jwtToken, _, err := jwtParser.ParseUnverified(token, claims)
	if err != nil {
		log.Printf("failed to parse JWT; %s", err.Error())
		os.Exit(1)
	}

	// FIXME-- use the unverified key material to lookup the signer's public key for verification below...
	if false {
		jwtToken, err = jwt.Parse(token, func(_jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := _jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("failed to resolve a valid JWT signing key; unsupported signing alg specified in header: %s", _jwtToken.Method.Alg())
			}

			var keypair *util.JWTKeypair

			var kid *string
			if kidhdr, ok := _jwtToken.Header["kid"].(string); ok {
				kid = &kidhdr
			}

			if kid != nil {
				keypair = jwtKeypairs[*kid]
			}

			if keypair == nil {
				for kid := range jwtKeypairs {
					keypair = jwtKeypairs[kid] // picks the last keypair...
				}
			}

			if keypair != nil {
				log.Printf("resolved keypair...")
				return &keypair.PublicKey, nil
			}

			return nil, errors.New("failed to resolve a valid JWT verification key")
		})
		if err != nil {
			log.Printf("failed to parse JWT; %s", err.Error())
			os.Exit(1)
		}

		log.Printf("%v", jwtToken)

		// claims, _ = jwtToken.Claims.(jwt.MapClaims)
		// if !claimsOk {
		// 	log.Printf("failed to parse claims in given bearer token")
		// 	os.Exit(1)
		// }
	}

	return claims
}

func requirePublicJWTVerifiers() {
	jwtKeypairs = map[string]*util.JWTKeypair{}

	keys, err := ident.GetJWKs()
	if err != nil {
		log.Printf("failed to resolve ident jwt keys; %s", err.Error())
	} else {
		for _, key := range keys {
			publicKey, err := pgputil.DecodeRSAPublicKeyFromPEM([]byte(key.PublicKey))
			if err != nil {
				log.Printf("failed to parse ident JWT public key; %s", err.Error())
			}

			sshPublicKey, err := ssh.NewPublicKey(publicKey)
			if err != nil {
				log.Printf("failed to resolve JWT public key fingerprint; %s", err.Error())
			}
			fingerprint := ssh.FingerprintLegacyMD5(sshPublicKey)

			jwtKeypairs[fingerprint] = &util.JWTKeypair{
				Fingerprint:  fingerprint,
				PublicKey:    *publicKey,
				PublicKeyPEM: &key.PublicKey,
			}

			log.Printf("ident jwt public key configured for verification; fingerprint: %s", fingerprint)
		}
	}
}

func init() {
	joinBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	joinBaselineWorkgroupCmd.MarkFlagRequired("organization")

	joinBaselineWorkgroupCmd.Flags().StringVar(&inviteJWT, "jwt", defaultNChainBaselineNetworkID, "JWT invitation token received from the inviting counterparty")
	joinBaselineWorkgroupCmd.MarkFlagRequired("jwt")
}
