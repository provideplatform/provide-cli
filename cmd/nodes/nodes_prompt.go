package nodes

import (
	"fmt"
	"strconv"

	"github.com/provideservices/provide-cli/cmd/common"
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
			common.Image = common.FreeInput("Image", "", "Mandatory")
		}
		if role == "" {
			role = common.FreeInput("Role", "", "Mandatory")
		}
		if optional {
			fmt.Println("Optional Flags:")
			if common.HealthCheckPath == "" {
				common.HealthCheckPath = common.FreeInput("Health Check Path", "", "")
			}
			if common.TCPIngressPorts == "" {
				common.TCPIngressPorts = common.FreeInput("TCP Ingress Ports", "", "")
			}
			if common.UDPIngressPorts == "" {
				common.UDPIngressPorts = common.FreeInput("UDP Ingress Ports", "", "")
			}
			if common.TaskRole == "" {
				common.TaskRole = common.FreeInput("Task Role", "", "")

			}
		}
		CreateNodeRun(cmd, args)
	case promptStepDelete:
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if common.NodeID == "" {
			common.NodeID = common.FreeInput("Node ID", "", "Mandatory")
		}
		deleteNodeRun(cmd, args)
	case promptStepLogs:
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if common.NodeID == "" {
			common.NodeID = common.FreeInput("Node ID", "", "Mandatory")
		}
		// Validation Number
		if page == 1 {
			result := common.FreeInput("Page", "1", "Number")
			page, _ = strconv.ParseUint(result, 10, 64)
		}
		// Validation Number
		if rpp == 100 {
			result := common.FreeInput("RPP", "100", "Number")
			rpp, _ = strconv.ParseUint(result, 10, 64)
		}
		nodeLogsRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
