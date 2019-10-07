package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var walletID string

var walletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Manage signing identities & cryptocurrency wallets",
	Long: `Various APIs are exposed to provide convenient access to elliptic-curve cryptography
(ECC) helper methods such as generating managed keypairs.

For convenience, it is also possible to generate decentralized keypairs with this utility. You
should securely store any keys generated using this API.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wallets unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(walletsCmd)
	walletsCmd.AddCommand(walletsListCmd)
	walletsCmd.AddCommand(walletsInitCmd)
}
