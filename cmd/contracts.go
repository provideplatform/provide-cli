package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Manage dapp smart contracts",
	Long:  `Smart contracts can be compiled from source using the CLI to streamline deployment and execution while enabling some advanced functionality such as resolving contract-internal opcodes to their respective `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("contracts unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(contractsCmd)
	contractsCmd.AddCommand(contractsListCmd)
	contractsCmd.AddCommand(contractsCompileCmd)
	contractsCmd.AddCommand(contractsDeployCmd)
}
