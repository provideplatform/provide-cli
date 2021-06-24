package participants

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var listBaselineWorkgroupParticipantsCmd = &cobra.Command{
	Use:   "list",
	Short: "List workgroup participants",
	Long:  `List the participating and invited parties in a baseline workgroup`,
	Run:   listParticipants,
}

func listParticipants(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listParticipantsRun(cmd *cobra.Command, args []string) {
	common.AuthorizeApplicationContext()
	common.AuthorizeOrganizationContext(false)

	participants, err := ident.ListApplicationOrganizations(common.OrganizationAccessToken, common.ApplicationID, map[string]interface{}{
		"type": "baseline",
	})
	if err != nil {
		log.Printf("failed to retrieve baseline workgroup participants; %s", err.Error())
		os.Exit(1)
	}

	invitations, err := ident.ListApplicationInvitations(common.ApplicationAccessToken, common.ApplicationID, map[string]interface{}{})
	if err != nil {
		// log.Printf("failed to retrieve invited baseline workgroup participants; %s", err.Error())
		// os.Exit(1)
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
	listBaselineWorkgroupParticipantsCmd.Flags().BoolVarP(&Optional, "Optional", "", false, "List all the Optional flags")
}
