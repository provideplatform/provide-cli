package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var connectorsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific connector",
	Long:  `Delete a specific connector by identifier and teardown any associated infrastructure`,
	Run:   deleteConnector,
}

func deleteConnector(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	status, _, err := provide.DeleteConnector(token, connectorID)
	if err != nil {
		log.Printf("Failed to delete connector with id: %s; %s", connectorID, err.Error())
		os.Exit(1)
	}
	if status != 204 {
		log.Printf("Failed to delete connector with id: %s; received status: %d", connectorID, status)
		os.Exit(1)
	}
	fmt.Printf("Deleted connector with id: %s", connectorID)
}

func init() {
	connectorsDeleteCmd.Flags().StringVar(&connectorID, "connector", "", "id of the connector")
	connectorsDeleteCmd.MarkFlagRequired("connector")

	connectorsDeleteCmd.Flags().StringVar(&applicationID, "application", "", "application id")
	connectorsDeleteCmd.MarkFlagRequired("application")
}
