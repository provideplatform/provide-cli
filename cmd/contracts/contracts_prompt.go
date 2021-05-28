package contracts

import (
	"strconv"

	"github.com/provideservices/provide-cli/cmd/common"
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
			contractExecMethod = common.FreeInput("Method", "", "Mandatory")
		}
		if common.ContractID == "" {
			common.ContractID = common.FreeInput("Contract ID", "", "Mandatory")
		}
		if optional {
			if common.AccountID == "" {
				common.RequireAccount(map[string]interface{}{})
			}
			if common.WalletID == "" {
				common.RequireWallet()
			}
			if contractExecValue == 0 {
				result := common.FreeInput("Value", "0", "")
				contractExecValue, _ = strconv.ParseUint(result, 10, 64)

			}
		}
		executeContract(cmd, args)
	case promptStepList:
		if optional {
			common.RequireApplication()
		}
		listContracts(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
