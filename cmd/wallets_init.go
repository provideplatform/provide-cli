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

var decentralized bool
var walletName string

var walletsInitCmd = &cobra.Command{
	Use:   "init [--decentralized|-d] [--network 024ff1ef-7369-4dee-969c-1918c6edb5d4]",
	Short: "Generate a new keypair for signing transactions and storing value",
	Long:  `Initialize a new wallet, which may be managed by Provide or you`,
	Run:   createWallet,
}

func createWallet(cmd *cobra.Command, args []string) {
	if decentralized {
		createDecentralizedWallet()
		return
	}

	createManagedWallet()
}

func createDecentralizedWallet() {
	publicKey, privateKey, err := provide.EVMGenerateKeyPair()
	if err != nil {
		log.Printf("Failed to genereate decentralized keypair; %s", err.Error())
		os.Exit(1)
	}
	secret := hex.EncodeToString(provide.FromECDSA(privateKey))
	keypairJSON, err := provide.EVMMarshalEncryptedKey(provide.HexToAddress(*publicKey), privateKey, secret)
	if err != nil {
		log.Printf("Failed to genereate decentralized keypair; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\n", *publicKey, string(keypairJSON))
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
		log.Printf("Failed to genereate keypair; %s", err.Error())
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
		fmt.Printf("Failed to generate keypair; %s", resp)
		os.Exit(1)
	}
}

func init() {
	walletsInitCmd.Flags().BoolVarP(&decentralized, "decentralized", "d", false, "if the generated keypair is decentralized")
	walletsInitCmd.Flags().StringVarP(&walletName, "name", "n", "", "human-readable name to associate withe the generated keypair")
	walletsInitCmd.MarkFlagRequired("network")
}
