/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package participants

import (
	participants_invitations "github.com/provideplatform/provide-cli/prvd/baseline/participants/invitations"
	participants_organizations "github.com/provideplatform/provide-cli/prvd/baseline/participants/organizations"
	participants_users "github.com/provideplatform/provide-cli/prvd/baseline/participants/users"
	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

const promptStepUsers = "Users"
const promptStepOrganizations = "Organizations"
const promptStepInvitations = "Invitations"

var emptyPromptArgs = []string{promptStepUsers, promptStepOrganizations, promptStepInvitations}
var emptyPromptLabel = "What would you like to do"

var Optional bool
var paginate bool

// var custodyPromptArgs = []string{"No", "Yes"}
// var custodyPromptLabel = "Would you like the participant to be a managed tenant?"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	// case promptStepInviteUser:
	// 	inviteUserRun(cmd, args)
	// case promptStepInviteOrganization:
	// 	inviteOrganizationRun(cmd, args)
	// case promptStepList:
	// 	if Optional {
	// 		fmt.Println("Optional Flags:")
	// 		if common.ApplicationID == "" {
	// 			common.RequireApplication()
	// 		}
	// 	}
	// 	page, rpp = common.PromptPagination(paginate, page, rpp)
	// 	listParticipantsRun(cmd, args)
	case promptStepUsers:
		participants_users.Optional = Optional
		participants_users.ParticipantsUsersCmd.Run(cmd, args)
	case promptStepOrganizations:
		participants_organizations.Optional = Optional
		participants_organizations.ParticipantsOrganizationsCmd.Run(cmd, args)
	case promptStepInvitations:
		participants_invitations.Optional = Optional
		participants_invitations.ParticipantsInvitationsCmd.Run(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
