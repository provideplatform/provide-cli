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

package applications

import (
	"fmt"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepDetails = "Details"
const promptStepInit = "Initialize"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepInit, promptStepList}
var emptyPromptLabel = "What would you like to do"

var axiomPromptArgs = []string{"Yes", "No"}
var axiomPromptLabel = "Would you like to make the application axiom compliant"

var accountPromptArgs = []string{"Yes", "No"}
var accountPromptLabel = "Would you like to make an account"

var walletPromptArgs = []string{"Yes", "No"}
var walletPromptLabel = "Would you like to set up a wallet"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		if applicationName == "" {
			applicationName = common.FreeInput("Application Name", "", common.MandatoryValidation)
		}
		if common.NetworkID == "" {
			common.RequireNetwork()
		}
		if optional {
			fmt.Println("Optional Flags:")
			if applicationType == "" {
				applicationType = common.FreeInput("Application Type", "", common.NoValidation)
			}
			if !axiom {
				result := common.SelectInput(axiomPromptArgs, axiomPromptLabel)
				axiom = result == "Yes"
			}
			if !withoutAccount {
				result := common.SelectInput(accountPromptArgs, accountPromptLabel)
				axiom = result == "Yes"
			}
			if !withoutWallet {
				result := common.SelectInput(walletPromptArgs, walletPromptLabel)
				axiom = result == "Yes"
			}
		}
		createApplication(cmd, args)
	case promptStepDetails:
		common.RequireApplication()
		fetchApplicationDetails(cmd, args)
	case promptStepList:
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listApplications(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
