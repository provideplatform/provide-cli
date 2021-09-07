package applications

import (
	"fmt"
	"strconv"

	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepDetails = "Details"
const promptStepInit = "Initialize"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepInit, promptStepList}
var emptyPromptLabel = "What would you like to do"

var baselinePromptArgs = []string{"Yes", "No"}
var baselinePromptLabel = "Would you like to make the application baseline compliant"

var accountPromptArgs = []string{"Yes", "No"}
var accountPromptLabel = "Would you like to make an account"

var walletPromptArgs = []string{"Yes", "No"}
var walletPromptLabel = "Would you like to set up a wallet"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		if applicationName == "" {
			applicationName = common.FreeInput("Application Name", "", common.MandatoryValidation)
		}
		if common.NetworkID == "" {
			common.RequireNetwork()
		}
		if optional {
			fmt.Println("Optional Flags:")
			if applicationType == "" {
				applicationType = common.FreeInput("Application Type", "", common.NoValidation)
			}
			if !baseline {
				result := common.SelectInput(baselinePromptArgs, baselinePromptLabel)
				baseline = result == "Yes"
			}
			if !withoutAccount {
				result := common.SelectInput(accountPromptArgs, accountPromptLabel)
				baseline = result == "Yes"
			}
			if !withoutWallet {
				result := common.SelectInput(walletPromptArgs, walletPromptLabel)
				baseline = result == "Yes"
			}
		}
		createApplication(cmd, args)
	case promptStepDetails:
		common.RequireApplication()
		fetchApplicationDetails(cmd, args)
	case promptStepList:
		if paginate {
			if page == common.DefaultPage {
				result := common.FreeInput("Page", fmt.Sprintf("%d", common.DefaultPage), common.MandatoryNumberValidation)
				page, _ = strconv.ParseUint(result, 10, 64)
			}
			if rpp == common.DefaultRpp {
				result := common.FreeInput("RPP", fmt.Sprintf("%d", common.DefaultRpp), common.MandatoryValidation)
				rpp, _ = strconv.ParseUint(result, 10, 64)
			}
		}
		listApplications(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
