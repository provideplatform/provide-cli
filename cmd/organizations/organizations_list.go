package organizations

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var organizationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of organizations",
	Long:  `Retrieve a list of organizations scoped to the authorized API token`,
	Run:   listOrganizations,
}

func listOrganizations(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	organizations, err := provide.ListOrganizations(token, params)
	if err != nil {
		log.Printf("Failed to retrieve organizations list; %s", err.Error())
		os.Exit(1)
	}
	for i := range organizations {
		organization := organizations[i]
		address := "0x"
		if addr, addrOk := organization.Metadata["address"].(string); addrOk {
			address = addr
		}
		result := fmt.Sprintf("%s\t%s\t%s\n", organization.ID.String(), *organization.Name, address)
		fmt.Print(result)
	}
}
