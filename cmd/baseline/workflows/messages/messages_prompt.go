package messages

import (
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepSend = "Send"

var items = map[string]string{"General Consistency": "general_consistency"}
var custodyPromptLabel = "Message Type"

var emptyPromptArgs = []string{promptStepSend}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepSend:
		if common.ApplicationID == "" {
			common.RequireWorkgroup()
		}
		if common.OrganizationID == "" {
			common.RequireOrganization()
		}
		if messageType == "" {
			opts := make([]string, 0)
			for k := range items {
				opts = append(opts, k)
			}
			messageType = common.SelectInput(opts, custodyPromptLabel)

		}
		if id == "" {
			id = common.FreeInput("ID", "", common.MandatoryValidation)
		}
		if baselineID == "" {
			baselineID = common.FreeInput("Baseline ID", "", common.NoValidation)
		}
		if data == "" {
			data = common.FreeInput("Data", "", common.JSONValidation)
		}
		sendMessageRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
