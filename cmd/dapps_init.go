package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var dappName string

var dappsInitCmd = &cobra.Command{
	Use:   "init --name 'my awesome dapp' --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new dapp",
	Long:  `Initialize a new dapp targeting a specified mainnet`,
	Run:   createApplication,
}

func createApplication(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{
		"name": dappName,
		"config": map[string]interface{}{
			"network_id": networkID,
		},
	}
	status, resp, err := provide.CreateApplication(token, params)
	if err != nil {
		log.Printf("Failed to initialize dapp; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		dapp := resp.(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", dapp["id"], dapp["name"])
		fmt.Print(result)
	}
}

func init() {
	dappsInitCmd.Flags().StringVar(&dappName, "name", "", "name of the dapp")
	dappsInitCmd.MarkFlagRequired("name")

	dappsInitCmd.Flags().StringVar(&networkID, "network", "", "network id (i.e., the network being used)")
	dappsInitCmd.MarkFlagRequired("network")
}
