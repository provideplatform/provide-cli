package connectors

import (
	"fmt"
	"strconv"

	"github.com/provideplatform/provide-cli/cmd/common"
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
		if paginate {
			if page == common.DefaultPage {
				result := common.FreeInput("Page", fmt.Sprintf("%d", common.DefaultPage), common.MandatoryNumberValidation)
				page, _ = strconv.ParseUint(result, 10, 64)
			}
			if rpp == common.DefaultRpp {
				result := common.FreeInput("RPP", fmt.Sprintf("%d", common.DefaultRpp), common.MandatoryValidation)
				rpp, _ = strconv.ParseUint(result, 10, 64)
			}
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
