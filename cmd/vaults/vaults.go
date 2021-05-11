package vaults

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VaultsCmd = &cobra.Command{
	Use:   "vaults",
	Short: "Manage vaults",
	Long: `Create and manage vaults and their associated keys and secrets.

Supports encrypt/decrypt and sign/verify operations for select key specs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Vaults command run")
		generalWalletPrompt(cmd, args, "empty")
	},
}

func init() {
	VaultsCmd.AddCommand(vaultsListCmd)
	VaultsCmd.AddCommand(vaultsInitCmd)
}
