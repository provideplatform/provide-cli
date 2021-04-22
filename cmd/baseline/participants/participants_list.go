package participants

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/provideservices/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var applicationAccessToken string
var organizationAccessToken string

var listBaselineWorkgroupParticipantsCmd = &cobra.Command{
	Use:   "list",
	Short: "List baseline workgroup participants",
	Long:  `List the participating and invited parties in a baseline workgroup`,
	Run:   listParticipants,
}

func authorizeApplicationContext() {
	token, err := ident.CreateToken(common.RequireUserAuthToken(), map[string]interface{}{
		"scope":          "offline_access",
		"application_id": common.ApplicationID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of application %s; %s", common.ApplicationID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		applicationAccessToken = *token.AccessToken
	}
}

func authorizeOrganizationContext() {
	token, err := ident.CreateToken(common.RequireUserAuthToken(), map[string]interface{}{
		"scope":           "offline_access",
		"organization_id": common.OrganizationID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of organization %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		organizationAccessToken = *token.AccessToken
	}
}

func listParticipants(cmd *cobra.Command, args []string) {
	authorizeApplicationContext()

	participants, err := ident.ListApplicationOrganizations(applicationAccessToken, common.ApplicationID, map[string]interface{}{
		"type": "baseline",
	})
	if err != nil {
		log.Printf("failed to retrieve baseline workgroup participants; %s", err.Error())
		os.Exit(1)
	}

	invitations, err := ident.ListApplicationInvitations(applicationAccessToken, common.ApplicationID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to retrieve invited baseline workgroup participants; %s", err.Error())
		os.Exit(1)
	}

	if len(participants) > 0 {
		fmt.Print("Organizations:\n")
	}

	for i := range participants {
		participant := participants[i]

		address := "0x"
		if addr, addrOk := participant.Metadata["address"].(string); addrOk {
			address = addr
		}

		var endpoint string
		if msgEndpoint, msgEndpointOk := participant.Metadata["messaging_endpoint"].(string); msgEndpointOk {
			endpoint = msgEndpoint
		}
		result := fmt.Sprintf("%s\t%s\t%s\t%s\n", participant.ID.String(), *participant.Name, address, endpoint)
		fmt.Print(result)
	}

	if len(invitations) > 0 {
		fmt.Print("\nPending Invitations:\n")
	}

	for i := range invitations {
		invitedParticipant := invitations[i]
		result := fmt.Sprintf("%s\t%s\n", invitedParticipant.ID.String(), invitedParticipant.Email)
		fmt.Print(result)
	}
}

func init() {
	listBaselineWorkgroupParticipantsCmd.Flags().StringVar(&common.ApplicationID, "workgroup", "", "workgroup identifier")
	listBaselineWorkgroupParticipantsCmd.MarkFlagRequired("workgroup")
}
