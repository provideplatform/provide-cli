package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dappsCmd = &cobra.Command{
	Use:   "dapps",
	Short: "Initialize and manage dapps and their associated resources",
	Long: `Initialized dapps are deployed to a specific mainnet and expose the following APIs:

	- API Tokens
	- Smart Contracts
	- Token Contracts
	- Wallets (also referred to as Signing Identities)
	- Oracles
	- Bridges
	- Connectors (i.e., IPFS)
	- Transactions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dapps unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(dappsCmd)
	dappsCmd.AddCommand(dappsListCmd)
	dappsCmd.AddCommand(dappsInitCmd)
}
