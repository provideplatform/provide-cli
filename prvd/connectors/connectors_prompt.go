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

package connectors

import (
	"strconv"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepInit = "Init"
const promptStepList = "List"
const promptStepDetails = "Details"
const promptStepDelete = "Delete"

var emptyPromptArgs = []string{promptStepInit, promptStepList, promptStepDetails, promptStepDelete}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		if connectorName == "" {
			connectorName = common.FreeInput("Connector Name", "", common.MandatoryValidation)
		}
		if connectorType == "" {
			connectorType = common.FreeInput("Connector Type", "", common.MandatoryValidation)
		}
		if common.ApplicationID == "" {
			common.RequireApplication()
		}
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if optional {
			if ipfsAPIPort == 5001 {
				result := common.FreeInput("IPFS API Port", "5001", common.NumberValidation)
				ipfsAPIPort, _ = strconv.ParseUint(result, 10, 64)
			}
			if ipfsGatewayPort == 8080 {
				result := common.FreeInput("IPFS Gateway Port", "8080", common.NumberValidation)
				ipfsGatewayPort, _ = strconv.ParseUint(result, 10, 64)
			}
		}
		createConnector(cmd, args)
	case promptStepList:
		if optional {
			common.RequireApplication()
		}
		page, rpp = common.PromptPagination(paginate, page, rpp)
		listConnectors(cmd, args)
	case promptStepDetails:
		common.RequireConnector(map[string]interface{}{})
		fetchConnectorDetails(cmd, args)
	case promptStepDelete:
		if common.ConnectorID == "" {
			common.RequireConnector(map[string]interface{}{})
		}
		if common.ApplicationID == "" {
			common.RequireApplication()
		}
		deleteConnector(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
