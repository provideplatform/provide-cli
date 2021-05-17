package organizations

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var organizationsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific organization",
	Long:  `Retrieve details for a specific organization by identifier, scoped to the authorized API token`,
	Run:   fetchOrganizationDetails,
}

func fetchOrganizationDetails(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepDetails)
}

func fetchOrganizationDetailsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		generalPrompt(cmd, args, "Details")
	}
	token := common.RequireUserAuthToken()
	params := map[string]interface{}{}
	organization, err := provide.GetOrganizationDetails(token, common.OrganizationID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for organization with id: %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve details for organization with id: %s; %s", common.OrganizationID, organization)
	// 	os.Exit(1)
	// }
	result := fmt.Sprintf("%s\t%s\n", organization.ID, *organization.Name)
	fmt.Print(result)
}

func init() {
	organizationsDetailsCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "id of the organization")
	// organizationsDetailsCmd.MarkFlagRequired("organization")
}
