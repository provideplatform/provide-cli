package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var walletID string

var walletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Generate managed and decentralized signing identities & cryptocurrency wallets",
	Long: `Various APIs are exposed to provide convenient access to
elliptic-curve cryptography (ECC) helper methods such as
generating managed keypairs.

It is also possible to generate decentralized keypairs. You
should securely store any decentralized keys generated using 
this API.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wallets unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(walletsCmd)
	walletsCmd.AddCommand(walletsListCmd)
	walletsCmd.AddCommand(walletsInitCmd)
}
