package users

import (
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptStepCreate = "Create"
var promptStepAuthenticate = "Authenticate"

var emptyPromptArgs = []string{promptStepCreate, promptStepAuthenticate}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepAuthenticate:
		authenticate(cmd, args)
	case promptStepCreate:
		create(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
