package baseledger

import (
	"github.com/provideplatform/provide-cli/cmd/baseledger/node"
	"github.com/provideplatform/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

const promptNode = "Node"

var emptyPromptArgs = []string{promptNode}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptNode:
		node.NodeCmd.Run(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
