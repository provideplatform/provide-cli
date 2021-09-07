package organizations

import (
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var organizationName string
var paginate bool

var organizationsInitCmd = &cobra.Command{
	Use:   "init --name 'Acme Inc.'",
	Short: "Initialize a new organization",
	Long:  `Initialize a new organization`,
	Run:   createOrganization,
}

func organizationConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"network_id": common.NetworkID,
	}

	return cfg
}
func createOrganization(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func createOrganizationRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"name":   organizationName,
		"config": organizationConfigFactory(),
	}
	organization, err := provide.CreateOrganization(token, params)
	if err != nil {
		log.Printf("Failed to initialize organization; %s", err.Error())
		os.Exit(1)
	}

	common.OrganizationID = organization.ID.String()
	log.Printf("initialized organization: %s\t%s\n", organizationName, common.OrganizationID)
}

func init() {
	organizationsInitCmd.Flags().StringVar(&organizationName, "name", "", "name of the organization")
	// organizationsInitCmd.MarkFlagRequired("name")
}
