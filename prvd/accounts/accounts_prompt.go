package accounts

import (
	"fmt"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepCustody = "Custody"
const promptStepInit = "Initialize"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepInit, promptStepList}
var emptyPromptLabel = "What would you like to do"

var accountTypePromptArgs = []string{"Managed", "Decentralised"}
var accountTypeLabel = "What type of Wallet would you like to create"

var custodyPromptArgs = []string{"No", "Yes"}
var custodyPromptLabel = "Would you like your wallet to be non-custodial?"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		common.SelectInput(accountTypePromptArgs, accountTypeLabel)
		generalPrompt(cmd, args, promptStepCustody)
	case promptStepCustody:
		if optional {
			fmt.Println("Optional Flags:")
			if !nonCustodial {
				nonCustodial = common.SelectInput(custodyPromptArgs, custodyPromptLabel) == "Yes"
			}
			if accountName == "" {
				accountName = common.FreeInput("Account Name", "", common.NoValidation)
			}
			if common.ApplicationID == "" {
				common.RequireApplication()
			}
			if common.OrganizationID == "" {
				common.RequireOrganization()
			}
		}
		CreateAccount(cmd, args)
	case promptStepList:
		if optional {
			fmt.Println("Optional Flags:")
			common.RequireApplication()
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listAccounts(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
