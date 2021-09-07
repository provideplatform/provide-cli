package api_tokens

import (
	"fmt"

	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepInit, promptStepList}
var emptyPromptLabel = "What would you like to do"

var refresTokenPromptArgs = []string{"Yes", "No"}
var refresTokenPromptLabel = "Would you like to set a refresh token"

var offlinePromptArgs = []string{"Yes", "No"}
var offlinePromptLabel = "Would you like to set offline access"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		if optional {
			if common.ApplicationID == "" {
				common.RequireApplication()
			}
			if common.OrganizationID == "" {
				common.RequireOrganization()
			}
			if !refreshToken {
				result := common.SelectInput(refresTokenPromptArgs, refresTokenPromptLabel)
				refreshToken = result == "Yes"
			}
			if !offlineAccess {
				result := common.SelectInput(offlinePromptArgs, offlinePromptLabel)
				offlineAccess = result == "Yes"
			}
			if refreshToken && offlineAccess {
				fmt.Println("⚠️  WARNING: You currently have both refresh and offline token set, Refresh token will take precedence")
			}
		}
		createAPIToken(cmd, args)
	case promptStepList:
		if optional {
			common.RequireApplication()
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listAPITokens(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
