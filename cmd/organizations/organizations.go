package organizations

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// var organization provide.Organization

var OrganizationsCmd = &cobra.Command{
	Use:   "organizations",
	Short: "Manage organizations",
	Long: `Create and manage organizations in the context of the following APIs:

	- Applications
	- Baseline Protocol
	- Tokens
	- Vaults`,
	Run: func(cmd *cobra.Command, args []string) {
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
	OrganizationsCmd.AddCommand(organizationsListCmd)
	OrganizationsCmd.AddCommand(organizationsInitCmd)
	OrganizationsCmd.AddCommand(organizationsDetailsCmd)
}
