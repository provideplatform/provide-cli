package users

import (
	"os"

	"github.com/spf13/cobra"
)

var UsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
	Long:  `Create and manage users and authenticate`,
	Run: func(cmd *cobra.Command, args []string) {
		authenticatePrompt(cmd, args)
		defer func() {
			if r := recover(); r != nil {
				os.Exit(1)
			}
		}()
	},
}

func init() {
	//no-op
}
