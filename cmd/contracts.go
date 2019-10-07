package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contract map[string]interface{}
var contractID string

var contractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Manage application smart contracts",
	Long:  `Compile smart contracts locally from source or execute previously-deployed contracts`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("contracts unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(contractsCmd)
	contractsCmd.AddCommand(contractsListCmd)
}
