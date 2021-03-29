package organizations

import (
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var organizationName string
var organizationType string

var organizationsInitCmd = &cobra.Command{
	Use:   "init --name 'Acme Inc.' --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new organization",
	Long:  `Initialize a new organization`,
	Run:   createOrganization,
}

func organizationConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"network_id": common.NetworkID,
	}

	if organizationType != "" {
		cfg["type"] = organizationType
	}

	return cfg
}

func createOrganization(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"name":   organizationName,
		"type":   organizationType,
		"config": organizationConfigFactory(),
	}
	organization, err := provide.CreateOrganization(token, params)
	if err != nil {
		log.Printf("Failed to initialize organization; %s", err.Error())
		os.Exit(1)
	}
	common.OrganizationID = organization.ID.String()
}

func init() {
	organizationsInitCmd.Flags().StringVar(&organizationName, "name", "", "name of the organization")
	organizationsInitCmd.MarkFlagRequired("name")
}
