package stack

import (
	"github.com/provideplatform/provide-cli/prvd/common"
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

var runBaselineStackCmd = &cobra.Command{
	Use:   "run",
	Short: "See `prvd baseline stack start --help` instead",
	Long: `Start a local baseline stack instance and connect to internal systems of record.

See: prvd baseline stack run --help instead. This command is deprecated and will be removed soon.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStackStart(cmd, args)
	},
}

func init() {
	StackCmd.AddCommand(logsBaselineStackCmd)
	StackCmd.AddCommand(runBaselineStackCmd)
	StackCmd.AddCommand(startBaselineStackCmd)
	StackCmd.AddCommand(stopBaselineStackCmd)
	StackCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the optional flags")
}
