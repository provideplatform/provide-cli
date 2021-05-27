package organizations

import (
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepDetails = "Details"
const promptStepInit = "Initialize"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepInit, promptStepList, promptStepDetails}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		organizationName = common.FreeInput("Organization Name", "", "Mandatory")
		createOrganizationRun(cmd, args)
	case promptStepList:
		listOrganizationsRun(cmd, args)
	case promptStepDetails:
		common.RequireOrganization()
		fetchOrganizationDetailsRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
