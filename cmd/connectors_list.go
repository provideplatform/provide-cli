package cmd

import (
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
	connectors = resp.([]interface{})
	for i := range connectors {
		connector := connectors[i].(map[string]interface{})
		config := connector["config"].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\t%s", connector["id"], connector["name"], connector["type"])
		if connector["type"] == connectorTypeIPFS {
			result = fmt.Sprintf("%s\t%s", result, config["api_url"])
		}
		fmt.Printf("%s\n", result)
	}
}

func init() {
	connectorsListCmd.Flags().StringVar(&applicationID, "application", "", "application identifier to filter connectors")
}
