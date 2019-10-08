package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const contractTypeRegistry = "registry"

var contract map[string]interface{}
var contracts []interface{}
var contractID string
var contractType string

var contractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Manage smart contracts",
	Long:  `Compile and deploy smart contracts locally from source or execute previously-deployed contracts`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("contracts unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(contractsCmd)
	contractsCmd.AddCommand(contractsListCmd)
	contractsCmd.AddCommand(contractsExecuteCmd)
}
