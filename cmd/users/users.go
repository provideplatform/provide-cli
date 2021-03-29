package users

import (
	"fmt"

	"github.com/spf13/cobra"
)

var UsersCmd = &cobra.Command{
	Use:   "organizations",
	Short: "Manage organizations",
	Long: `Create and manage organizations in the context of the following APIs:

	- Applications
	- Baseline Protocol
	- Tokens
	- Vaults`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("organizations unimplemented")
	},
}

func init() {
	UsersCmd.AddCommand(authenticateCmd)
}
