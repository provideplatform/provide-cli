package keys

import (
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys",
	Long:  `Create and manage cryptographic keys. Supports symmetric and asymmetric key specs with encrypt/decrypt and sign/verify operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)
		generalPrompt(cmd, args, "")
	},
}

func init() {
	KeysCmd.AddCommand(keysListCmd)
	KeysCmd.AddCommand(keysInitCmd)
	KeysCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
