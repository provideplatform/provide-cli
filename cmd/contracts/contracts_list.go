package contracts

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var pagination *common.Pagination

var contractsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of contracts",
	Long:  `Retrieve a list of contracts scoped to the authorized API token`,
	Run:   listContracts,
}

func listContracts(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	contracts, resp, err := provide.ListContracts(token, params)
	if err != nil {
		log.Printf("Failed to retrieve contracts list; %s", err.Error())
		os.Exit(1)
	}
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(contracts))
	for i := range contracts {
		contract := contracts[i]
		result := fmt.Sprintf("%s\t%s\t%s\n", contract.ID.String(), *contract.Address, *contract.Name)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listContracts,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	contractsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter contracts")
	contractsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	contractsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	contractsListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	contractsListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of contracts to retrieve per page")
}
