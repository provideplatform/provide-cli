package users

import (
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var UsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
	Long:  `Create and manage users and authenticate`,
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
	UsersCmd.AddCommand(createCmd)
	UsersCmd.AddCommand(showIDCmd)
}
