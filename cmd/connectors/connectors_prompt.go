package connectors

import (
	"strconv"

	"github.com/provideservices/provide-cli/cmd/common"
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
			connectorName = common.FreeInput("Connector Name")
		}
		if connectorType == "" {
			connectorType = common.FreeInput("Connector Type")
		}
		if common.ApplicationID == "" {
			common.RequireApplication()
		}
		if common.NetworkID == "" {
			common.RequirePublicNetwork()
		}
		if optional {
			// Validation number
			if ipfsAPIPort == 5001 {
				result := common.FreeInput("IPFS API Port")
				ipfsAPIPort, _ = strconv.ParseUint(result, 10, 64)
			}
			// Validation number
			if ipfsGatewayPort == 8080 {
				result := common.FreeInput("IPFS Gateway Port")
				ipfsGatewayPort, _ = strconv.ParseUint(result, 10, 64)
			}
		}
		createConnector(cmd, args)
	case promptStepList:
		if optional {
			common.RequireApplication()
		}
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
