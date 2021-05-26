package workgroups

import (
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/api_tokens"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/provideservices/provide-go/api/baseline"
	"github.com/spf13/cobra"
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

	api_tokens.RequirePublicJWTVerifiers() // FIXME...
	claims := parseJWT(inviteJWT)
	log.Printf("resolved baseline claims containing invitation for workgroup: %s", *claims.Baseline.WorkgroupID)

	common.ApplicationID = *claims.Baseline.WorkgroupID
	common.AuthorizeOrganizationContext(true)
	authorizeApplicationContext()

	// initWorkgroupContract()

	common.RequireOrganizationVault()
	requireOrganizationKeys()
	common.RegisterWorkgroupOrganization(common.ApplicationID)
	// common.RequireOrganizationEndpoints(nil)

	configureBaselineStack(inviteJWT, claims)
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
		jwtToken, err = api_tokens.ParseJWT(token)
		if err != nil {
			log.Printf("failed to parse JWT; %s", err.Error())
			os.Exit(1)
		}

		log.Printf("%v", jwtToken)
	}

	return claims
}

// configureBaselineStack initializes a workgroup in the context of the running baseline stack
func configureBaselineStack(jwt string, claims *InviteClaims) {
	token := common.RequireAPIToken()
	_, err := baseline.CreateWorkgroup(token, map[string]interface{}{
		"token": jwt,
	})
	if err != nil {
		// log.Printf("failed to configure baseline stack to support joined workgroup; %s", err.Error())
		os.Exit(1)
	}
	log.Printf("configured baseline workgroup on local stack: %s", *claims.Baseline.WorkgroupID)
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

func init() {
	joinBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	joinBaselineWorkgroupCmd.Flags().StringVar(&inviteJWT, "jwt", "", "JWT invitation token received from the inviting counterparty")
	joinBaselineWorkgroupCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")

}
