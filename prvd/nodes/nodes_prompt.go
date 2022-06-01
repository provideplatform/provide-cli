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

package nodes

import (
	"fmt"
	"strconv"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepLogs = "Logs"
const promptStepInit = "Initialize"
const promptStepDelete = "Delete"

var emptyPromptArgs = []string{promptStepInit, promptStepLogs, promptStepDelete}
var emptyPromptLabel = "What would you like to do"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if common.Image == "" {
			common.Image = common.FreeInput("Image", "", common.MandatoryValidation)
		}
		if role == "" {
			role = common.FreeInput("Role", "", common.MandatoryValidation)
		}
		if optional {
			fmt.Println("Optional Flags:")
			if common.HealthCheckPath == "" {
				common.HealthCheckPath = common.FreeInput("Health Check Path", "", common.NoValidation)
			}
			if common.TCPIngressPorts == "" {
				common.TCPIngressPorts = common.FreeInput("TCP Ingress Ports", "", common.NoValidation)
			}
			if common.UDPIngressPorts == "" {
				common.UDPIngressPorts = common.FreeInput("UDP Ingress Ports", "", common.NoValidation)
			}
			if common.TaskRole == "" {
				common.TaskRole = common.FreeInput("Task Role", "", common.NoValidation)

			}
		}
		CreateNodeRun(cmd, args)
	case promptStepDelete:
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if common.NodeID == "" {
			common.NodeID = common.FreeInput("Node ID", "", common.MandatoryValidation)
		}
		deleteNodeRun(cmd, args)
	case promptStepLogs:
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if common.NodeID == "" {
			common.NodeID = common.FreeInput("Node ID", "", common.MandatoryValidation)
		}
		// Validation Number
		if page == 1 {
			result := common.FreeInput("Page", "1", common.MandatoryNumberValidation)
			page, _ = strconv.ParseUint(result, 10, 64)
		}
		// Validation Number
		if rpp == 100 {
			result := common.FreeInput("RPP", "100", common.MandatoryValidation)
			rpp, _ = strconv.ParseUint(result, 10, 64)
		}
		nodeLogsRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
