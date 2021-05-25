package baseline

import (
	"fmt"

	"github.com/provideservices/provide-cli/cmd/baseline/stack"
	"github.com/provideservices/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

const promptStackCustody = "Stack"
const promptWorkgroupsInit = "Workgroups"
const promptWorkflowsList = "Workflows"
const promptParticipantList = "Participants"

var emptyPromptArgs = []string{promptStackCustody, promptWorkgroupsInit, promptWorkflowsList, promptParticipantList}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStackCustody:
		fmt.Print(optional)
		stack.StackCmd.Run(cmd, args)
		stack.OptionalStack = optional
	case promptWorkgroupsInit:
	case promptWorkflowsList:
	case promptParticipantList:
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		stack.OptionalStack = optional
		generalPrompt(cmd, args, result)
	}
}
