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

var nonCustodial bool
var walletName string

var walletsInitCmd = &cobra.Command{
	Use:   "init [--non-custodial|-nc]",
	Short: "Generate a new HD wallet for deterministically managing accounts, signing transactions and storing value",
	Long:  `Initialize a new HD wallet, which may be managed by Provide or you`,
	Run:   createWallet,
}

func createWallet(cmd *cobra.Command, args []string) {
	if nonCustodial {
		createDecentralizedWallet()
		return
	}

	createManagedWallet()
}

func createDecentralizedWallet() {
	publicKey, privateKey, err := provide.EVMGenerateKeyPair()
	if err != nil {
		log.Printf("Failed to genereate decentralized HD wallet; %s", err.Error())
		os.Exit(1)
	}
	secret := hex.EncodeToString(provide.FromECDSA(privateKey))
	walletJSON, err := provide.EVMMarshalEncryptedKey(provide.HexToAddress(*publicKey), privateKey, secret)
	if err != nil {
		log.Printf("Failed to genereate decentralized HD wallet; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\n", *publicKey, string(walletJSON))
	fmt.Print(result)
}

func createManagedWallet() {
	token := requireAPIToken()
	params := map[string]interface{}{
		"network_id": networkID,
	}
	if walletName != "" {
		params["name"] = walletName
	}
	status, resp, err := provide.CreateWallet(token, params)
	if err != nil {
		log.Printf("Failed to genereate HD wallet; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		wallet := resp.(map[string]interface{})
		walletID = wallet["id"].(string)
		result := fmt.Sprintf("Wallet %s\t%s\n", wallet["id"], wallet["address"])
		if name, nameOk := wallet["name"].(string); nameOk {
			result = fmt.Sprintf("Wallet %s\t%s - %s\n", name, wallet["id"], wallet["address"])
		}
		appWalletKey := buildConfigKeyWithApp(walletConfigKeyPartial, applicationID)
		if !viper.IsSet(appWalletKey) {
			viper.Set(appWalletKey, wallet["id"])
			viper.WriteConfig()
		}
		fmt.Print(result)
	} else {
		fmt.Printf("Failed to generate HD wallet; %s", resp)
		os.Exit(1)
	}
}

func init() {
	walletsInitCmd.Flags().BoolVarP(&nonCustodial, "non-custodial", "", false, "if the generated HD wallet is custodial")
	walletsInitCmd.Flags().StringVarP(&walletName, "name", "n", "", "human-readable name to associate with the generated HD wallet")
}
