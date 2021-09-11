package accounts

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"
	"github.com/spf13/cobra"
)

var pagination *common.Pagination

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of signing identities",
	Long:  `Retrieve a list of signing identities (accounts) scoped to the authorized API token`,
	Run:   listAccounts,
}

func listAccounts(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	results, resp, err := provide.ListAccounts(token, params)
	if err != nil {
		log.Printf("Failed to retrieve accounts list; %s", err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve accounts list; received status: %d", status)
	// 	os.Exit(1)
	// }
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(results))
	for i := range results {
		account := results[i]
		result := fmt.Sprintf("%s\t%s\n", account.ID.String(), account.Address)
		// TODO-- when account.Name exists... result = fmt.Sprintf("%s\t%s - %s\n", name, account, *account.Address)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listAccounts,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	accountsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter accounts")
	accountsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	accountsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	accountsListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	accountsListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of accounts to retrieve per page")
}
