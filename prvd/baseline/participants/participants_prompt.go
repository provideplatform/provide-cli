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
	"fmt"
	"strconv"

	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

const promptStepInvite = "Invite"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepInvite, promptStepList}
var emptyPromptLabel = "What would you like to do"

var custodyPromptArgs = []string{"No", "Yes"}
var custodyPromptLabel = "Would you like the participant to be a managed tenant?"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInvite:
		if Optional {
			fmt.Println("Optional Flags:")
			if common.ApplicationID == "" {
				common.RequireApplication()
			}
			if common.OrganizationID == "" {
				common.RequireOrganization()
			}
			if !managedTenant {
				managedTenant = common.SelectInput(custodyPromptArgs, custodyPromptLabel) == "Yes"
			}
			if name == "" {
				name = common.FreeInput("Wallet Name", "", common.NoValidation)
			}
			if email == "" {
				email = common.FreeInput("Wallet Purpose", "", common.NoValidation)
			}
			if permissions == 0 {
				permissions, _ = strconv.Atoi(common.FreeInput("Wallet Purpose", "0", common.NoValidation))
			}
		}
		inviteParticipantRun(cmd, args)
	case promptStepList:
		if Optional {
			fmt.Println("Optional Flags:")
			if common.ApplicationID == "" {
				common.RequireApplication()
			}
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listParticipantsRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
