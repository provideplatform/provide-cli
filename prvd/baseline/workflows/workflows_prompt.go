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

package workflows

import (
	"github.com/provideplatform/provide-cli/prvd/baseline/workflows/messages"
	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepList = "List"
const promptStepDetails = "Details"

const promptStepMessages = "Messages"

var emptyPromptArgs = []string{promptStepInit, promptStepMessages}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		initWorkflowRun(cmd, args)
	case promptStepList:
		listWorkflowsRun(cmd, args)
	case promptStepDetails:
		fetchWorkflowDetailsRun(cmd, args)
	case promptStepMessages:
		messages.Optional = Optional
		messages.MessagesCmd.Run(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
