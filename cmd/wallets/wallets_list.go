package wallets

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var pagination *common.Pagination

var walletsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of custodial HD wallets",
	Long:  `Retrieve a list of HD wallets scoped to the authorized API token`,
	Run:   listWallets,
}

func listWallets(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listWalletsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	results, resp, err := provide.ListWallets(token, params)
	if err != nil {
		log.Printf("Failed to retrieve wallets list; %s", err.Error())
		os.Exit(1)
	}
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(results))
	for i := range results {
		wallet := results[i]
		result := fmt.Sprintf("%s\t%s\n", wallet.ID.String(), *wallet.PublicKey)
		// FIXME-- when wallet.Name exists... result = fmt.Sprintf("Wallet %s\t%s - %s\n", wallet.Name, wallet.ID.String(), *wallet.Address)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listWalletsRun,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	walletsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter HD wallets")
	walletsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	walletsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	walletsListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	walletsListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of wallets to retrieve per page")
}
