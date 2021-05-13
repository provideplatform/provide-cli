package connectors

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var connectorsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve details for a specific connector",
	Long:  `Retrieve details for a specific connector by identifier, scoped to the authorized API token`,
	Run:   fetchConnectorDetails,
}

func fetchConnectorDetails(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, "Details")

	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	connector, err := provide.GetConnectorDetails(token, common.ConnectorID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for connector with id: %s; %s", common.ConnectorID, err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve details for connector with id: %s; received status: %d", common.ConnectorID, status)
	// 	os.Exit(1)
	// }
	var config map[string]interface{}
	json.Unmarshal(*connector.Config, &config)
	result := fmt.Sprintf("%s\t%s\t%s", connector.ID.String(), *connector.Name, *connector.Type)
	if *connector.Type == connectorTypeIPFS {
		result = fmt.Sprintf("%s\t%s", result, config["api_url"])
	}
	fmt.Printf("%s\n", result)
}

func init() {
	connectorsDetailsCmd.Flags().StringVar(&common.ConnectorID, "connector", "", "id of the connector")
	// connectorsDetailsCmd.MarkFlagRequired("connector")
}
