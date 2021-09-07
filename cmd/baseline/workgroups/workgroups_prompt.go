package workgroups

import (
	"fmt"

	"github.com/provideplatform/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepList = "List"
const promptStepJoin = "Join"

var emptyPromptArgs = []string{promptStepInit, promptStepList, promptStepJoin}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		if Optional {
			fmt.Println("Optional Flags:")
			if common.NetworkID == "" {
				common.RequirePublicNetwork()
			}
			if common.OrganizationID == "" {
				common.RequireOrganization()
			}
			if common.MessagingEndpoint == "" {
				common.MessagingEndpoint = common.FreeInput("Messaging Endpoint", "", common.NoValidation)
			}
			if name == "" {
				name = common.FreeInput("Name", "", common.NoValidation)
			}
		}
		initWorkgroupRun(cmd, args)
	case promptStepList:
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listWorkgroupsRun(cmd, args)
	case promptStepJoin:
		if Optional {
			fmt.Println("Optional Flags:")
			if common.OrganizationID == "" {
				common.RequireOrganization()
			}
			if inviteJWT == "" {
				inviteJWT = common.FreeInput("JWT Invite", "", common.NoValidation)
			}
		}
		joinWorkgroupRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
