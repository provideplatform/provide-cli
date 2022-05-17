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

package networks

import (
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepList = "List"
const promptStepDisable = "Disable"

var emptyPromptArgs = []string{promptStepInit, promptStepList, promptStepDisable}
var emptyPromptLabel = "What would you like to do"

var publicPromptArgs = []string{"Yes", "No"}
var publicPromptLabel = "Would you like the network to be public"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		// Validation non-null
		if chain == "" {
			chain = common.FreeInput("Chain", "", common.NoValidation)
		}
		if nativeCurrency == "" {
			nativeCurrency = common.FreeInput("Native Currency", "", common.NoValidation)
		}
		if platform == "" {
			platform = common.FreeInput("Platform", "", common.NoValidation)
		}
		if protocolID == "" {
			protocolID = common.FreeInput("Protocol ID", "", common.NoValidation)
		}
		if networkName == "" {
			networkName = common.FreeInput("Network Name", "", common.NoValidation)
		}
		CreateNetwork(cmd, args)
	case promptStepList:
		if optional {
			result := common.SelectInput(publicPromptArgs, publicPromptLabel)
			public = result == "Yes"
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listNetworks(cmd, args)
	case promptStepDisable:
		common.RequireNetwork()
		disableNetwork(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
