package organizations

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var pagination *common.Pagination

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
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	}
	organizations, resp, err := provide.ListOrganizations(token, params)
	if err != nil {
		log.Printf("Failed to retrieve organizations list; %s", err.Error())
		os.Exit(1)
	}
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(organizations))
	for i := range organizations {
		organization := organizations[i]
		address := "0x"
		if addr, addrOk := organization.Metadata["address"].(string); addrOk {
			address = addr
		}
		result := fmt.Sprintf("%s\t%s\t%s\n", organization.ID.String(), *organization.Name, address)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listOrganizationsRun,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	organizationsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	organizationsListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	organizationsListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of organizations to retrieve per page")
}
