package wallets

import (
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var WalletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Manage HD wallets and accounts",
	Long: `Generate hierarchical deterministic (HD) wallets and sign transactions.

More documentation forthcoming.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)
		generalPrompt(cmd, args, "")

		defer func() {
			if r := recover(); r != nil {
				os.Exit(1)
			}
		}()
	},
}

func init() {
	WalletsCmd.AddCommand(walletsListCmd)
	WalletsCmd.AddCommand(walletsInitCmd)
	WalletsCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
