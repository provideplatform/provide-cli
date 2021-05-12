package wallets

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var WalletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Manage HD wallets and accounts",
	Long: `Generate hierarchical deterministic (HD) wallets and sign transactions.

More documentation forthcoming.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Wallet command run")
		generalPrompt(cmd, args, "")

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Prompt Exit\n")
				os.Exit(1)
			}
		}()
	},
}

func init() {
	WalletsCmd.AddCommand(walletsListCmd)
	WalletsCmd.AddCommand(walletsInitCmd)
}
