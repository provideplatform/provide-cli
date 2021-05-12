package vaults

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var VaultsCmd = &cobra.Command{
	Use:   "vaults",
	Short: "Manage vaults",
	Long: `Create and manage vaults and their associated keys and secrets.

Supports encrypt/decrypt and sign/verify operations for select key specs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Vaults command run")
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
	VaultsCmd.AddCommand(vaultsListCmd)
	VaultsCmd.AddCommand(vaultsInitCmd)
}
