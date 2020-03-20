package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var networksDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a specific network",
	Long:  `Disable a specific network by identifier`,
	Run:   disableNetwork,
}

func disableNetwork(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	status, _, err := provide.UpdateNetwork(token, networkID, map[string]interface{}{
		"enabled": false,
	})
	if err != nil {
		log.Printf("Failed to disable network with id: %s; %s", networkID, err.Error())
		os.Exit(1)
	}
	if status != 204 {
		log.Printf("Failed to disable network with id: %s; received status: %d", networkID, status)
		os.Exit(1)
	}
	fmt.Printf("Disabled network with id: %s", networkID)
}

func init() {
	networksDisableCmd.Flags().StringVar(&networkID, "network", "", "id of the network")
	networksDisableCmd.MarkFlagRequired("network")
}
