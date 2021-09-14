package networks

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var public bool

var pagination *common.Pagination

var networksListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of networks",
	Long:  `Retrieve a list of networks scoped to the authorized API token`,
	Run:   listNetworks,
}

func listNetworks(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	}
	if public {
		params["public"] = "true"
	}
	networks, resp, err := provide.ListNetworks(token, params)
	if err != nil {
		log.Printf("Failed to retrieve networks list; %s", err.Error())
		os.Exit(1)
	}
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(networks))
	for i := range networks {
		network := networks[i]
		result := fmt.Sprintf("%s\t%s\n", network.ID.String(), *network.Name)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listNetworks,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	networksListCmd.Flags().BoolVarP(&public, "public", "p", false, "filter private networks (false by default)")
	networksListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	networksListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	networksListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	networksListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of networks to retrieve per page")
}
