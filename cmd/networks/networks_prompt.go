package networks

import (
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepList = "List"
const promptStepDisable = "Disable"

var emptyPromptArgs = []string{promptStepInit, promptStepList, promptStepDisable}
var emptyPromptLabel = "What would you like to do"

var publicPromptArgs = []string{"Yes", "No"}
var publicPromptLabel = "Would you like the network to be public"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		// Validation non-null
		if chain == "" {
			chain = common.FreeInput("Chain", "", common.NoValidation)
		}
		if nativeCurrency == "" {
			nativeCurrency = common.FreeInput("Native Currency", "", common.NoValidation)
		}
		if platform == "" {
			platform = common.FreeInput("Platform", "", common.NoValidation)
		}
		if protocolID == "" {
			protocolID = common.FreeInput("Protocol ID", "", common.NoValidation)
		}
		if networkName == "" {
			networkName = common.FreeInput("Network Name", "", common.NoValidation)
		}
		CreateNetwork(cmd, args)
	case promptStepList:
		if optional {
			result := common.SelectInput(publicPromptArgs, publicPromptLabel)
			public = result == "Yes"
		}
		listNetworks(cmd, args)
	case promptStepDisable:
		common.RequireNetwork()
		disableNetwork(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
