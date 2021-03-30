package users

import (
	"fmt"

	"github.com/spf13/cobra"
)

var UsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
	Long:  `Create and manage users and authenticate`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users unimplemented")
	},
}

func init() {
	// no-op
}
