package wallets

import (
	"fmt"
	"strconv"

	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

const promptStepCustody = "Custody"
const promptStepInit = "Initialize"
const promptStepList = "List"

var custodyPromptArgs = []string{"No", "Yes"}
var custodyPromptLabel = "Would you like your wallet to be non-custodial?"

var walletTypePromptArgs = []string{"Managed", "Decentralised"}
var walletTypeLabel = "What type of Wallet would you like to create"

var emptyPromptArgs = []string{promptStepInit, promptStepList}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		common.SelectInput(walletTypePromptArgs, walletTypeLabel)
		generalPrompt(cmd, args, promptStepCustody)
	case promptStepCustody:
		if optional {
			fmt.Println("Optional Flags:")
			if !nonCustodial {
				nonCustodial = common.SelectInput(custodyPromptArgs, custodyPromptLabel) == "Yes"
			}
			if walletName == "" {
				walletName = common.FreeInput("Wallet Name", "", common.NoValidation)
			}
			if purpose == 44 {
				purpose, _ = strconv.Atoi(common.FreeInput("Wallet Purpose", "44", common.NumberValidation))
			}
		}
		CreateWalletRun(cmd, args)
	case promptStepList:
		if optional {
			fmt.Println("Optional Flags:")
			if common.ApplicationID == "" {
				common.RequireApplication()
			}
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listWalletsRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
