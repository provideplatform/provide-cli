package applications

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var pagination *common.Pagination

var applicationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of applications",
	Long:  `Retrieve a list of applications scoped to the authorized API token`,
	Run:   listApplications,
}

func listApplications(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	}
	applications, resp, err := provide.ListApplications(token, params)
	if err != nil {
		log.Printf("Failed to retrieve applications list; %s", err.Error())
		os.Exit(1)
	}
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(applications))
	for i := range applications {
		application := applications[i]
		result := fmt.Sprintf("%s\t%s\n", application.ID.String(), *application.Name)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listApplications,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	applicationsListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	applicationsListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of applications to retrieve per page")
	applicationsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
