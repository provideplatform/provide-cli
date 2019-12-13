package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountName string

var accountsInitCmd = &cobra.Command{
	Use:   "init [--non-custodial|-nc] [--network 024ff1ef-7369-4dee-969c-1918c6edb5d4]",
	Short: "Generate a new keypair for signing transactions and storing value",
	Long:  `Initialize a new account, which may be managed by Provide or you`,
	Run:   createAccount,
}

func createAccount(cmd *cobra.Command, args []string) {
	if nonCustodial {
		createDecentralizedAccount()
		return
	}

	createManagedAccount()
}

func createDecentralizedAccount() {
	publicKey, privateKey, err := provide.EVMGenerateKeyPair()
	if err != nil {
		log.Printf("Failed to genereate nonCustodial keypair; %s", err.Error())
		os.Exit(1)
	}
	secret := hex.EncodeToString(provide.FromECDSA(privateKey))
	keypairJSON, err := provide.EVMMarshalEncryptedKey(provide.HexToAddress(*publicKey), privateKey, secret)
	if err != nil {
		log.Printf("Failed to genereate nonCustodial keypair; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\n", *publicKey, string(keypairJSON))
	fmt.Print(result)
}

func createManagedAccount() {
	token := requireAPIToken()
	params := map[string]interface{}{
		"network_id": networkID,
	}
	if accountName != "" {
		params["name"] = accountName
	}
	status, resp, err := provide.CreateAccount(token, params)
	if err != nil {
		log.Printf("Failed to genereate keypair; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		account := resp.(map[string]interface{})
		accountID = account["id"].(string)
		result := fmt.Sprintf("Account %s\t%s\n", account["id"], account["address"])
		if name, nameOk := account["name"].(string); nameOk {
			result = fmt.Sprintf("Account %s\t%s - %s\n", name, account["id"], account["address"])
		}
		appAccountKey := buildConfigKeyWithApp(accountConfigKeyPartial, applicationID)
		if !viper.IsSet(appAccountKey) {
			viper.Set(appAccountKey, account["id"])
			viper.WriteConfig()
		}
		fmt.Print(result)
	} else {
		fmt.Printf("Failed to generate keypair; %s", resp)
		os.Exit(1)
	}
}

func init() {
	accountsInitCmd.Flags().BoolVarP(&nonCustodial, "non-custodial", "", false, "if the generated keypair is non-custodial")
	accountsInitCmd.Flags().StringVarP(&accountName, "name", "n", "", "human-readable name to associate with the generated keypair")
	accountsInitCmd.MarkFlagRequired("network")
}
