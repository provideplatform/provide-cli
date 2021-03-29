package contracts

import (
	"fmt"

	"github.com/spf13/cobra"
)

const contractTypeRegistry = "registry"

var contract map[string]interface{}
var contracts []interface{}
var contractType string

var ContractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Manage smart contracts",
	Long:  `Compile and deploy smart contracts locally from source or execute previously-deployed contracts`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("contracts unimplemented")
	},
}

func init() {
	ContractsCmd.AddCommand(contractsListCmd)
	ContractsCmd.AddCommand(contractsExecuteCmd)
}
