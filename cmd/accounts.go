package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var accountID string

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage signing identities & accounts",
	Long: `Various APIs are exposed to provide convenient access to elliptic-curve cryptography
(ECC) helper methods such as generating managed (custodial) keypairs.

For convenience, it is also possible to generate keypairs with this utility which you (or your application)
is then responsible for securing. You should securely store any keys generated using this API. If you are
looking for hierarchical deterministic support, check out the wallets API.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("accounts unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsInitCmd)
}
