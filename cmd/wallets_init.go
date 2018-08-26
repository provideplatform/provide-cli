package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var decentralized bool
var walletName string

var walletsInitCmd = &cobra.Command{
	Use:   "init [--decentralized|-d]",
	Short: "Generate a new keypair for signing transactions and storing value",
	Long:  `Initialize a new wallet, which may be managed or decentralized.`,
	Run:   createWallet,
}

func createWallet(cmd *cobra.Command, args []string) {
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
