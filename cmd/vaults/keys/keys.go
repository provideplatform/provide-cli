package keys

import (
	"os"

	"github.com/spf13/cobra"
)

var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys",
	Long: `Create and manage cryptographic keys.

Supports symmetric and asymmetric key specs with encrypt/decrypt and sign/verify operations.

Docs: https://docs.provide.services/vault/api-reference/keys`,
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
	KeysCmd.AddCommand(keysListCmd)
	KeysCmd.AddCommand(keysInitCmd)
}
