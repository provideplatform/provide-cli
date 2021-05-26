package applications

import (
	"os"

	"github.com/spf13/cobra"
)

const applicationTypeMessageBus = "message_bus"

var application map[string]interface{}

var ApplicationsCmd = &cobra.Command{
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
		generalPrompt(cmd, args, "")

		defer func() {
			if r := recover(); r != nil {
				os.Exit(1)
			}
		}()
	},
}

func init() {
	ApplicationsCmd.AddCommand(applicationsListCmd)
	ApplicationsCmd.AddCommand(applicationsInitCmd)
	ApplicationsCmd.AddCommand(applicationsDetailsCmd)
	ApplicationsCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
