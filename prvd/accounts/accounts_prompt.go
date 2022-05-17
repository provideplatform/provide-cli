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

package accounts

import (
	"fmt"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepCustody = "Custody"
const promptStepInit = "Initialize"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepInit, promptStepList}
var emptyPromptLabel = "What would you like to do"

var accountTypePromptArgs = []string{"Managed", "Decentralised"}
var accountTypeLabel = "What type of Wallet would you like to create"

var custodyPromptArgs = []string{"No", "Yes"}
var custodyPromptLabel = "Would you like your wallet to be non-custodial?"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		common.SelectInput(accountTypePromptArgs, accountTypeLabel)
		generalPrompt(cmd, args, promptStepCustody)
	case promptStepCustody:
		if optional {
			fmt.Println("Optional Flags:")
			if !nonCustodial {
				nonCustodial = common.SelectInput(custodyPromptArgs, custodyPromptLabel) == "Yes"
			}
			if accountName == "" {
				accountName = common.FreeInput("Account Name", "", common.NoValidation)
			}
			if common.ApplicationID == "" {
				common.RequireApplication()
			}
			if common.OrganizationID == "" {
				common.RequireOrganization()
			}
		}
		CreateAccount(cmd, args)
	case promptStepList:
		if optional {
			fmt.Println("Optional Flags:")
			common.RequireApplication()
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listAccounts(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
