package vaults

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-cli/prvd/vaults/keys"
)

var VaultsCmd = &cobra.Command{
	Use:   "vaults",
	Short: "Manage vaults",
	Long: `Create and manage vaults and their associated keys and secrets.

Vaults support select symmetric and asymmetric key specs for encrypt/decrypt and sign/verify operations.

Docs: https://docs.provide.services/vault`,
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
	VaultsCmd.AddCommand(vaultsListCmd)
	VaultsCmd.AddCommand(vaultsInitCmd)

	VaultsCmd.AddCommand(keys.KeysCmd)
	VaultsCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	VaultsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
