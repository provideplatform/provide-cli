package node

import (
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var NodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Interact with a local baseledger node",
	Long:  `Create, manage and interact with a local baseledger node.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

func init() {
	NodeCmd.AddCommand(startBaseledgerNodeCmd)
	NodeCmd.AddCommand(stopBaseledgerNodeCmd)
}
