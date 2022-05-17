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

package contracts

import (
	"strconv"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepExecute = "Execute"
const promptStepList = "List"

var emptyPromptArgs = []string{promptStepExecute, promptStepList}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepExecute:
		if contractExecMethod == "" {
			contractExecMethod = common.FreeInput("Method", "", common.MandatoryValidation)
		}
		if common.ContractID == "" {
			common.ContractID = common.FreeInput("Contract ID", "", common.MandatoryValidation)
		}
		if optional {
			if common.AccountID == "" {
				common.RequireAccount(map[string]interface{}{})
			}
			if common.WalletID == "" {
				common.RequireWallet()
			}
			if contractExecValue == 0 {
				result := common.FreeInput("Value", "0", common.NumberValidation)
				contractExecValue, _ = strconv.ParseUint(result, 10, 64)

			}
		}
		executeContract(cmd, args)
	case promptStepList:
		if optional {
			common.RequireApplication()
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
	case "":
		listContracts(cmd, args)
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
