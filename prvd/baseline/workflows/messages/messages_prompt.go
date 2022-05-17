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

package messages

import (
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepSend = "Send"

var items = map[string]string{"General Consistency": "general_consistency"}
var custodyPromptLabel = "Message Type"

var emptyPromptArgs = []string{promptStepSend}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepSend:
		if common.ApplicationID == "" {
			common.RequireWorkgroup()
		}
		if common.OrganizationID == "" {
			common.RequireOrganization()
		}
		if messageType == "" {
			opts := make([]string, 0)
			for k := range items {
				opts = append(opts, k)
			}
			value := common.SelectInput(opts, custodyPromptLabel)
			messageType = items[value]
		}
		if id == "" {
			id = common.FreeInput("ID", "", common.MandatoryValidation)
		}
		if baselineID == "" {
			baselineID = common.FreeInput("Baseline ID", "", common.NoValidation)
		}
		if data == "" {
			data = common.FreeInput("Data", "", common.JSONValidation)
		}
		sendMessageRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
