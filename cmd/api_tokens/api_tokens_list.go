package api_tokens

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var pagination *common.Pagination

var apiTokensListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of API tokens",
	Long:  `Retrieve a list of API tokens scoped to the authorized API token`,
	Run:   listAPITokens,
}

func listAPITokens(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	results, resp, err := provide.ListTokens(token, params)
	if err != nil {
		log.Printf("Failed to retrieve API tokens list; %s", err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve API tokens list; received status: %d", status)
	// 	os.Exit(1)
	// }
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(results))
	for i := range results {
		apiToken := results[i]
		result := fmt.Sprintf("%s\t%s\n", apiToken.ID.String(), *apiToken.Token)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listAPITokens,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	apiTokensListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter API tokens")
	apiTokensListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	apiTokensListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	apiTokensListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	apiTokensListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of API tokens to retrieve per page")
}
