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

package workgroups

import (
	"fmt"

	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepList = "List"
const promptStepDetails = "Details"
const promptStepJoin = "Join"

var emptyPromptArgs = []string{promptStepInit, promptStepList, promptStepJoin}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepList:
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listWorkgroupsRun(cmd, args)
	case promptStepDetails:
		fetchWorkgroupDetailsRun(cmd, args)
	case promptStepInit:
		if Optional {
			fmt.Println("Optional Flags:")
			if common.NetworkID == "" {
				common.RequireL1Network()
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
