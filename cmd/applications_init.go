package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var applicationName string
var applicationType string
var withoutAPIToken bool
var withoutAccount bool
var withoutWallet bool

var applicationsInitCmd = &cobra.Command{
	Use:   "init --name 'my app' --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new application",
	Long:  `Initialize a new application targeting a specified mainnet`,
	Run:   createApplication,
}

func applicationConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"network_id": networkID,
	}

	if applicationType != "" {
		cfg["type"] = applicationType
	}

	return cfg
}

func createApplication(cmd *cobra.Command, args []string) {
	if withoutAPIToken && !withoutWallet {
		fmt.Println("Cannot create an application that has a wallet but no API token.")
		os.Exit(1)
	}
	token := requireAPIToken()
	params := map[string]interface{}{
		"name":   applicationName,
		"type":   applicationType,
		"config": applicationConfigFactory(),
	}
	status, resp, err := provide.CreateApplication(token, params)
	if err != nil {
		log.Printf("Failed to initialize application; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		response := resp.(map[string]interface{})
		application = response["application"].(map[string]interface{})
		token := response["token"].(map[string]interface{})
		applicationID = application["id"].(string)
		applicationToken := token["token"].(string)

		appAPITokenKey := buildConfigKeyWithApp(apiTokenConfigKeyPartial, applicationID)
		if !viper.IsSet(appAPITokenKey) {
			viper.Set(appAPITokenKey, applicationToken)
			viper.WriteConfig()
		}
		fmt.Printf("Application API Token\t%s\n", applicationToken)

		result := fmt.Sprintf("%s\t%s\n", application["id"], application["name"])
		fmt.Print(result)
	}
	if !withoutAccount {
		createAccount(cmd, args)
	}
	if !withoutWallet {
		createWallet(cmd, args)
	}
}

func init() {
	applicationsInitCmd.Flags().StringVar(&applicationName, "name", "", "name of the application")
	applicationsInitCmd.MarkFlagRequired("name")

	applicationsInitCmd.Flags().StringVar(&applicationType, "type", "", "application type (i.e., message_bus)")

	applicationsInitCmd.Flags().StringVar(&networkID, "network", "", "target network id")
	applicationsInitCmd.MarkFlagRequired("network")

	applicationsInitCmd.Flags().BoolVar(&withoutWallet, "without-account", false, "do not create a new account (signing identity) for this application")
	applicationsInitCmd.Flags().BoolVar(&withoutWallet, "without-wallet", false, "do not create a new HD wallet for this application")
}
