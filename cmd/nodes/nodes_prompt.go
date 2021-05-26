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
		// Validation non-null
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if common.Image == "" {
			common.Image = common.FreeInput("Image")
		}
		if role == "" {
			role = common.FreeInput("Role")
		}
		if optional {
			fmt.Println("Optional Flags:")
			if common.HealthCheckPath == "" {
				common.HealthCheckPath = common.FreeInput("Health Check Path")
			}
			if common.TCPIngressPorts == "" {
				common.TCPIngressPorts = common.FreeInput("TCP Ingress Ports")
			}
			if common.UDPIngressPorts == "" {
				common.UDPIngressPorts = common.FreeInput("UDP Ingress Ports")
			}
			if common.TaskRole == "" {
				common.TaskRole = common.FreeInput("Task Role")

			}
		}
		CreateNodeRun(cmd, args)
	case promptStepDelete:
		// Validation non-null
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if common.NodeID == "" {
			common.NodeID = common.FreeInput("Node ID")
		}
		deleteNodeRun(cmd, args)
	case promptStepLogs:
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if common.NodeID == "" {
			common.NodeID = common.FreeInput("Node ID")
		}
		// Validation Number
		if page == 1 {
			result := common.FreeInput("Page")
			page, _ = strconv.ParseUint(result, 10, 64)
		}
		// Validation Number
		if rpp == 100 {
			result := common.FreeInput("RPP")
			rpp, _ = strconv.ParseUint(result, 10, 64)
		}
		nodeLogsRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
