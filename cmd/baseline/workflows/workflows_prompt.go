package workflows

import (
	"github.com/provideservices/provide-cli/cmd/baseline/workflows/messages"
	"github.com/provideservices/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepMessages = "Messages"

var emptyPromptArgs = []string{promptStepInit, promptStepMessages}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		if workgroupID == "" {
			common.RequireWorkgroup()
		}
		if name == "" {
			name = common.FreeInput("Name")
		}
		initWorkflowRun(cmd, args)
	case promptStepMessages:
		messages.Optional = Optional
		messages.MessagesCmd.Run(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
