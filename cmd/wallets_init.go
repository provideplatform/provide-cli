package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var decentralized bool
var walletName string

var walletsInitCmd = &cobra.Command{
	Use:   "init [--decentralized|-d] [--network 024ff1ef-7369-4dee-969c-1918c6edb5d4]",
	Short: "Generate a new keypair for signing transactions and storing value",
	Long:  `Initialize a new wallet, which may be managed or decentralized.`,
	Run:   createWallet,
}

func createWallet(cmd *cobra.Command, args []string) {
	if decentralized {
		publicKey, privateKey, err := provide.GenerateKeyPair()
		if err != nil {
			log.Printf("Failed to genereate decentralized keypair; %s", err.Error())
			os.Exit(1)
		}
		secret := hex.EncodeToString(crypto.FromECDSA(privateKey))
		keypairJSON, err := provide.MarshalEncryptedKey(common.HexToAddress(*publicKey), privateKey, secret)
		if err != nil {
			log.Printf("Failed to genereate decentralized keypair; %s", err.Error())
			os.Exit(1)
		}
		result := fmt.Sprintf("%s\t%s\n", *publicKey, string(keypairJSON))
		fmt.Print(result)
		return
	}

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
		result := fmt.Sprintf("%s\t%s\n", wallet["id"], wallet["address"])
		if name, nameOk := wallet["name"].(string); nameOk {
			result = fmt.Sprintf("%s\t%s - %s\n", wallet["id"], name, wallet["address"])
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
