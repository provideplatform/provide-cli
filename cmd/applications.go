package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var application map[string]interface{}
var applicationID string

var applicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Manage applications",
	Long: `Create and manage logical applications which target a specific network and expose the following APIs:

	- API Tokens
	- Smart Contracts
	- Token Contracts
	- Signing Identities (wallets)
	- Oracles
	- Bridges
	- Connectors (i.e., IPFS)
	- Payment Hubs
	- Transactions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("applications unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(applicationsCmd)
	applicationsCmd.AddCommand(applicationsListCmd)
	applicationsCmd.AddCommand(applicationsInitCmd)
}
