package wallets

import (
	"fmt"

	"github.com/spf13/cobra"
)

var WalletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Manage HD wallets and accounts",
	Long: `Generate hierarchical deterministic (HD) wallets and sign transactions.

More documentation forthcoming.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Wallet command run")
		generalWalletPrompt(cmd, args, "empty")
	},
}

func init() {
	WalletsCmd.AddCommand(walletsListCmd)
	WalletsCmd.AddCommand(walletsInitCmd)
}
