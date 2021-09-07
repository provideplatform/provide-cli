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

var page uint64
var rpp uint64

var connectorsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of connectors",
	Long:  `Retrieve a list of connectors scoped to the authorized API token`,
	Run:   listConnectors,
}

func listConnectors(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	connectors, err := provide.ListConnectors(token, params)
	if err != nil {
		log.Printf("Failed to retrieve connectors list; %s", err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve connectors list; received status: %d", status)
	// 	os.Exit(1)
	// }
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
}

func init() {
	connectorsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter connectors")
	connectorsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	connectorsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	connectorsListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	connectorsListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of connectors to retrieve per page")
}
