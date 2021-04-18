package baseline

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/provideservices/provide-cli/cmd/baseline/stack"
	"github.com/provideservices/provide-cli/cmd/baseline/workflows"
	"github.com/provideservices/provide-cli/cmd/baseline/workgroups"
)

var BaselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Interact with the baseline protocol",
	Long:  `Interact with the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("baseline unimplemented")
	},
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "See `prvd baseline stack --help` instead",
	Long: `Create, manage and interact with local baseline stack instances.

See: prvd baseline stack --help instead. This command is deprecated and will be removed soon.`,
	Run: func(cmd *cobra.Command, args []string) {
		stack.StackCmd.Run(cmd, args)
	},
}

func init() {
	BaselineCmd.AddCommand(proxyCmd)
	BaselineCmd.AddCommand(stack.StackCmd)
	BaselineCmd.AddCommand(workgroups.WorkgroupsCmd)
	BaselineCmd.AddCommand(workflows.WorkflowsCmd)
}
