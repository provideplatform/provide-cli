package stack

import (
	"fmt"

	"github.com/spf13/cobra"
)

var StackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Interact with a local baseline stack",
	Long:  `Create, manage and interact with local baseline stack instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stack unimplemented")
	},
}

func init() {
	StackCmd.AddCommand(logsBaselineStackCmd)
	StackCmd.AddCommand(runBaselineStackCmd)
	StackCmd.AddCommand(stopBaselineStackCmd)
}
