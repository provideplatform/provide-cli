package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var applicationName string
var withoutAPIToken bool
var withoutWallet bool

var applicationsInitCmd = &cobra.Command{
	Use:   "init --name 'my app' --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new application",
	Long:  `Initialize a new application targeting a specified mainnet`,
	Run:   createApplication,
}

func createApplication(cmd *cobra.Command, args []string) {
	if withoutAPIToken && !withoutWallet {
		fmt.Println("Cannot create an application that has a wallet but no API token.")
		os.Exit(1)
	}
	token := requireAPIToken()
	params := map[string]interface{}{
		"name": applicationName,
		"config": map[string]interface{}{
			"network_id": networkID,
		},
	}
	status, resp, err := provide.CreateApplication(token, params)
	if err != nil {
		log.Printf("Failed to initialize application; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		application := resp.(map[string]interface{})
		applicationID = application["id"].(string)
		result := fmt.Sprintf("%s\t%s\n", application["name"], application["id"])
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
	applicationsInitCmd.Flags().StringVar(&applicationName, "name", "", "name of the application")
	applicationsInitCmd.MarkFlagRequired("name")

	applicationsInitCmd.Flags().StringVar(&networkID, "network", "", "target network id")
	applicationsInitCmd.MarkFlagRequired("network")

	applicationsInitCmd.Flags().BoolVar(&withoutAPIToken, "without-api-token", false, "do not create a new API token for this application")
	applicationsInitCmd.Flags().BoolVar(&withoutWallet, "without-wallet", false, "do not create a new wallet (signing identity) for this application")
}
