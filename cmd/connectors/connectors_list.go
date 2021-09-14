package connectors

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var pagination *common.Pagination

var connectorsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of connectors",
	Long:  `Retrieve a list of connectors scoped to the authorized API token`,
	Run:   listConnectors,
}

func listConnectors(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	connectors, resp, err := provide.ListConnectors(token, params)
	if err != nil {
		log.Printf("Failed to retrieve connectors list; %s", err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve connectors list; received status: %d", status)
	// 	os.Exit(1)
	// }
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(connectors))
	for i := range connectors {
		connector := connectors[i]
		var config map[string]interface{}
		json.Unmarshal(*connector.Config, &config)
		result := fmt.Sprintf("%s\t%s\t%s", connector.ID.String(), *connector.Name, *connector.Type)
		if *connector.Type == connectorTypeIPFS {
			result = fmt.Sprintf("%s\t%s", result, config["api_url"])
		}
		fmt.Printf("%s\n", result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listConnectors,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	connectorsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter connectors")
	connectorsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	connectorsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	connectorsListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	connectorsListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of connectors to retrieve per page")
}
