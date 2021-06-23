package stack

import (
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var StackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Interact with a local baseline stack",
	Long:  `Create, manage and interact with local baseline stack instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

func init() {
	StackCmd.AddCommand(logsBaselineStackCmd)
	StackCmd.AddCommand(runBaselineStackCmd)
	StackCmd.AddCommand(stopBaselineStackCmd)
	StackCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the optional flags")
}
