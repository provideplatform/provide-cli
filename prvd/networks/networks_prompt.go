package networks

import (
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepList = "List"
const promptStepDisable = "Disable"

var emptyPromptArgs = []string{promptStepList, promptStepDisable}
var emptyPromptLabel = "What would you like to do"

var publicPromptArgs = []string{"Yes", "No"}
var publicPromptLabel = "Would you like the network to be public"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepList:
		if optional {
			result := common.SelectInput(publicPromptArgs, publicPromptLabel)
			public = result == "Yes"
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listNetworks(cmd, args)
	case promptStepDisable:
		common.RequireNetwork()
		disableNetwork(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
