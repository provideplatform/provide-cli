package organizations

import (
	"fmt"

	"github.com/spf13/cobra"

	provide "github.com/provideservices/provide-go/api/ident"
)

var organization provide.Organization

var OrganizationsCmd = &cobra.Command{
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
	OrganizationsCmd.AddCommand(organizationsListCmd)
	OrganizationsCmd.AddCommand(organizationsInitCmd)
	OrganizationsCmd.AddCommand(organizationsDetailsCmd)
}
