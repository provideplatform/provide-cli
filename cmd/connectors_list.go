package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var connectorsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of connectors",
	Long:  `Retrieve a list of connectors scoped to the authorized API token`,
	Run:   listConnectors,
}

func listConnectors(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	if applicationID != "" {
		params["application_id"] = applicationID
	}
	status, resp, err := provide.ListConnectors(token, params)
	if err != nil {
		log.Printf("Failed to retrieve connectors list; %s", err.Error())
		os.Exit(1)
	}
	if status != 200 {
		log.Printf("Failed to retrieve connectors list; received status: %d", status)
		os.Exit(1)
	}
	for i := range resp.([]interface{}) {
		connector := resp.([]interface{})[i].(map[string]interface{})
		config := connector["config"].(map[string]interface{})
		connectorConfigJSON, _ := json.Marshal(config)
		result := fmt.Sprintf("%s\t%s\t%s\n", connector["id"], connector["type"], connectorConfigJSON)
		fmt.Print(result)
	}
}

func init() {
	connectorsListCmd.Flags().StringVar(&applicationID, "application", "", "application identifier to filter connectors")
}
