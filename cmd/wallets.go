package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var walletID string

var walletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Manage HD wallets and accounts",
	Long: `Generate hierarchical deterministic (HD) wallets and sign transactions.

More documentation forthcoming.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wallets unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(walletsCmd)
	walletsCmd.AddCommand(walletsListCmd)
	walletsCmd.AddCommand(walletsInitCmd)
}
