package vaults

import (
	"fmt"

	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepInit, promptStepList}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		if name == "" {
			name = common.FreeInput("Vault Name", "", common.NoValidation)
		}
		if optional {
			fmt.Println("Optional Flags:")
			if description == "" {
				description = common.FreeInput("Vault Description", "", common.NoValidation)
			}
			if common.ApplicationID == "" {
				common.RequireApplication()
			}
			if common.OrganizationID == "" {
				common.RequireOrganization()
			}
		}
		createVaultRun(cmd, args)
	case promptStepList:
		if optional {
			fmt.Println("Optional Flags:")
			if common.ApplicationID == "" {
				common.RequireApplication()
			}
			if common.OrganizationID == "" {
				common.RequireOrganization()
			}
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listVaultsRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
