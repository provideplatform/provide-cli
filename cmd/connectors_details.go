package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var connectorsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve details for a specific connector",
	Long:  `Retrieve details for a specific connector by identifier, scoped to the authorized API token`,
	Run:   fetchConnectorDetails,
}

func fetchConnectorDetails(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	status, resp, err := provide.GetConnectorDetails(token, connectorID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for connector with id: %s; %s", connectorID, err.Error())
		os.Exit(1)
	}
	if status != 200 {
		log.Printf("Failed to retrieve details for connector with id: %s; received status: %d", connectorID, status)
		os.Exit(1)
	}
	connector := resp.(map[string]interface{})
	config := connector["config"].(map[string]interface{})
	result := fmt.Sprintf("%s\t%s\t%s", connector["id"], connector["name"], connector["type"])
	if connector["type"] == connectorTypeIPFS {
		result = fmt.Sprintf("%s\t%s", result, config["api_url"])
	}
	fmt.Printf("%s\n", result)
}

func init() {
	connectorsDetailsCmd.Flags().StringVar(&connectorID, "connector", "", "id of the connector")
	connectorsDetailsCmd.MarkFlagRequired("connector")
}
