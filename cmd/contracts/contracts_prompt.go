package contracts

import (
	"fmt"
	"strconv"

	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepExecute = "Execute"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepExecute, promptStepList}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepExecute:
		if contractExecMethod == "" {
			contractExecMethod = common.FreeInput("Method", "", common.MandatoryValidation)
		}
		if common.ContractID == "" {
			common.ContractID = common.FreeInput("Contract ID", "", common.MandatoryValidation)
		}
		if optional {
			if common.AccountID == "" {
				common.RequireAccount(map[string]interface{}{})
			}
			if common.WalletID == "" {
				common.RequireWallet()
			}
			if contractExecValue == 0 {
				result := common.FreeInput("Value", "0", common.NumberValidation)
				contractExecValue, _ = strconv.ParseUint(result, 10, 64)

			}
		}
		executeContract(cmd, args)
	case promptStepList:
		if optional {
			common.RequireApplication()
		}
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
	case "":
		listContracts(cmd, args)
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
