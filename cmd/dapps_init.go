package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var dappName string
var withoutAPIToken bool
var withoutWallet bool

var dappsInitCmd = &cobra.Command{
	Use:   "init --name 'my awesome dapp' --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new dapp",
	Long:  `Initialize a new dapp targeting a specified mainnet`,
	Run:   createApplication,
}

func createApplication(cmd *cobra.Command, args []string) {
	if withoutAPIToken && !withoutWallet {
		fmt.Println("Cannot create an application that has a wallet but no API token.")
		os.Exit(1)
	}
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
		applicationID = dapp["id"].(string)
		result := fmt.Sprintf("%s\t%s\n", dapp["name"], dapp["id"])
		fmt.Print(result)
	}
	if !withoutAPIToken {
		createAPIToken(cmd, args)
	}
	if !withoutWallet {
		createWallet(cmd, args)
	}
}

func init() {
	dappsInitCmd.Flags().StringVar(&dappName, "name", "", "name of the dapp")
	dappsInitCmd.MarkFlagRequired("name")

	dappsInitCmd.Flags().StringVar(&networkID, "network", "", "network id (i.e., the dapp mainnet)")
	dappsInitCmd.MarkFlagRequired("network")

	dappsInitCmd.Flags().BoolVar(&withoutAPIToken, "without-api-token", false, "do not create a new API token for this dapp")
	dappsInitCmd.Flags().BoolVar(&withoutWallet, "without-wallet", false, "do not create a new wallet (signing identity) for this dapp")
}
