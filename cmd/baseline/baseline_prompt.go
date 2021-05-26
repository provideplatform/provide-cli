package baseline

import (
	"github.com/provideservices/provide-cli/cmd/baseline/participants"
	"github.com/provideservices/provide-cli/cmd/baseline/stack"
	"github.com/provideservices/provide-cli/cmd/baseline/workflows"
	"github.com/provideservices/provide-cli/cmd/baseline/workgroups"

	"github.com/provideservices/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

const promptStack = "Stack"
const promptWorkgroups = "Workgroups"
const promptWorkflows = "Workflows"
const promptParticipant = "Participants"

var emptyPromptArgs = []string{promptStack, promptWorkgroups, promptWorkflows, promptParticipant}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStack:
		stack.Optional = Optional
		stack.StackCmd.Run(cmd, args)
	case promptWorkgroups:
		workgroups.Optional = Optional
		workgroups.WorkgroupsCmd.Run(cmd, args)
	case promptWorkflows:
		workflows.Optional = Optional
		workflows.WorkflowsCmd.Run(cmd, args)
	case promptParticipant:
		participants.Optional = Optional
		participants.ParticipantsCmd.Run(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
