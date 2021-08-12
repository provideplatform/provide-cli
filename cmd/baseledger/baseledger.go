package baseledger

import (
	"github.com/spf13/cobra"

	"github.com/provideplatform/provide-cli/cmd/baseledger/node"
	"github.com/provideplatform/provide-cli/cmd/common"
)

var Optional bool

var BaseledgerCmd = &cobra.Command{
	Use:   "baseledger",
	Short: "Manage a local baseledger node",
	Long:  `Manage a local baseledger node.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

func init() {
	BaseledgerCmd.AddCommand(node.NodeCmd)
}
