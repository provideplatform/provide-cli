package contracts

import (
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
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
	ContractsCmd.AddCommand(contractsListCmd)
	ContractsCmd.AddCommand(contractsExecuteCmd)
	ContractsCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")

}
