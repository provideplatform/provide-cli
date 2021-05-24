package vaults

import (
	"os"

	"github.com/spf13/cobra"
)

var VaultsCmd = &cobra.Command{
	Use:   "vaults",
	Short: "Manage vaults",
	Long: `Create and manage vaults and their associated keys and secrets.

Supports encrypt/decrypt and sign/verify operations for select key specs.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	VaultsCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
