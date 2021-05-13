package connectors

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var connectorsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific connector",
	Long:  `Delete a specific connector by identifier and teardown any associated infrastructure`,
	Run:   deleteConnector,
}

func deleteConnector(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, "Delete")

	token := common.RequireAPIToken()
	err := provide.DeleteConnector(token, common.ConnectorID)
	if err != nil {
		log.Printf("Failed to delete connector with id: %s; %s", common.ConnectorID, err.Error())
		os.Exit(1)
	}
	// if status != 204 {
	// 	log.Printf("Failed to delete connector with id: %s; received status: %d", common.ConnectorID, status)
	// 	os.Exit(1)
	// }
	fmt.Printf("Deleted connector with id: %s", common.ConnectorID)
}

func init() {
	connectorsDeleteCmd.Flags().StringVar(&common.ConnectorID, "connector", "", "id of the connector")
	// connectorsDeleteCmd.MarkFlagRequired("connector")

	connectorsDeleteCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application id")
	// connectorsDeleteCmd.MarkFlagRequired("application")
}
