package workflows

import (
	"github.com/provideservices/provide-cli/cmd/baseline/workflows/messages"
	"github.com/spf13/cobra"
)

var WorkflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "Interact with a baseline workflows",
	Long:  `Create, manage and interact with workflows via the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		generalPrompt(cmd, args, "")
	},
}

func init() {
	WorkflowsCmd.AddCommand(initBaselineWorkflowCmd)
	WorkflowsCmd.AddCommand(messages.MessagesCmd)
	WorkflowsCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")

}
