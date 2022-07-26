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

package baseline

import (
	"github.com/provideplatform/provide-cli/prvd/baseline/participants"
	"github.com/provideplatform/provide-cli/prvd/baseline/stack"
	"github.com/provideplatform/provide-cli/prvd/baseline/subject_accounts"
	"github.com/provideplatform/provide-cli/prvd/baseline/workflows"
	"github.com/provideplatform/provide-cli/prvd/baseline/workgroups"

	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

const promptStack = "Stack"
const promptWorkgroups = "Workgroups"
const promptWorkflows = "Workflows"
const promptParticipant = "Participants"
const promptSubjectAccounts = "Subject-Accounts"

var emptyPromptArgs = []string{promptStack, promptWorkgroups, promptWorkflows, promptParticipant, promptSubjectAccounts}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStack:
		stack.Optional = Optional
		stack.StackCmd.Run(cmd, args)
	case promptWorkgroups:
		workgroups.Optional = Optional
		workgroups.WorkgroupsCmd.Run(cmd, args)
	case promptWorkflows:
		workflows.Optional = Optional
		workflows.WorkflowsCmd.Run(cmd, args)
	case promptParticipant:
		participants.Optional = Optional
		participants.ParticipantsCmd.Run(cmd, args)
	case promptSubjectAccounts:
		subject_accounts.Optional = Optional
		subject_accounts.SubjectAccountsCmd.Run(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
