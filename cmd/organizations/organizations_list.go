package organizations

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var organizationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of organizations",
	Long:  `Retrieve a list of organizations scoped to the authorized API token`,
	Run:   listOrganizations,
}

func listOrganizations(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listOrganizationsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
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

func init() {
	organizationsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	organizationsListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	organizationsListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of organizations to retrieve per page")
}
