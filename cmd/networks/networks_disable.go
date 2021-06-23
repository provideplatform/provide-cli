package networks

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var networksDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a specific network",
	Long:  `Disable a specific network by identifier`,
	Run:   disableNetwork,
}

func disableNetwork(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	err := provide.UpdateNetwork(token, common.NetworkID, map[string]interface{}{
		"enabled": false,
	})
	if err != nil {
		log.Printf("Failed to disable network with id: %s; %s", common.NetworkID, err.Error())
		os.Exit(1)
	}
	// if status != 204 {
	// 	log.Printf("Failed to disable network with id: %s; received status: %d", common.NetworkID, status)
	// 	os.Exit(1)
	// }
	fmt.Printf("Disabled network with id: %s", common.NetworkID)
}

func init() {
	networksDisableCmd.Flags().StringVar(&common.NetworkID, "network", "", "id of the network")
	// networksDisableCmd.MarkFlagRequired("network")
}
