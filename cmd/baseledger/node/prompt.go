package node

import (
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepStart = "Start"
const promptStepStop = "Stop"

var emptyPromptArgs = []string{promptStepStart, promptStepStop}
var emptyPromptLabel = "What would you like to do"

// General prompt
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepStart:
		startBaseledgerNode(cmd, args)
	case promptStepStop:
		stopBaseledgerNode(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
