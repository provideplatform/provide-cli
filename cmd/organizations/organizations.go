package organizations

import (
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
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
	OrganizationsCmd.AddCommand(organizationsListCmd)
	OrganizationsCmd.AddCommand(organizationsInitCmd)
	OrganizationsCmd.AddCommand(organizationsDetailsCmd)
}
