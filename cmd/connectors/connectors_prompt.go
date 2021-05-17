package connectors

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepInit = "Init"
const promptStepList = "List"
const promptStepDetails = "Details"
const promptStepDelete = "Delete"
const promptStepSummary = "Summary"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		mandatoryInitFlags()
		if flagPrompt(cmd, args) {
			optionalInitFlags(cmd, args)
		}
	case promptStepList:
		if flagPrompt(cmd, args) {
			applicationIDFlagPrompt()
		}
		summary(cmd, args, promptArgs)
	case promptStepDetails:
		mandatoryDetailFlags()
		summary(cmd, args, promptArgs)
	case promptStepDelete:
		mandatoryDeleteFlags()
		summary(cmd, args, promptArgs)
	case promptStepSummary:
		summary(cmd, args, promptArgs)
	case "":
		emptyPrompt(cmd, args)
	default:
		fmt.Println("no-ops")
	}
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepInit, promptStepList, promptStepDetails, promptStepDelete},
	}

	_, result, _ := prompt.Run()

	promptArgs = append(promptArgs, result)

	generalPrompt(cmd, args, result)
}

func flagPrompt(cmd *cobra.Command, args []string) bool {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"No", "Yes"},
	}

	_, flagResult, err := flagPrompt.Run()

	if err != nil {
		os.Exit(1)
		return false
	}

	if flagResult == "Yes" {
		return true
	} else {
		generalPrompt(cmd, args, promptStepSummary)
		return false
	}
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == promptStepInit {
		createConnector(cmd, args)
	}
	if promptArgs[0] == promptStepList {
		listConnectors(cmd, args)
	}
	if promptArgs[0] == promptStepDetails {
		fetchConnectorDetails(cmd, args)
	}
	if promptArgs[0] == promptStepDelete {
		deleteConnector(cmd, args)
	}
}

func mandatoryInitFlags() {
	if connectorName == "" {
		connectorNameFlagPrompt()
	}
	if connectorType == "" {
		connectorTypeFlagPrompt()
	}
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.NetworkID == "" {
		networkIDFlagPrompt()
	}
}

func mandatoryDeleteFlags() {
	if common.ConnectorID == "" {
		connectorIDFlagPrompt()
	}
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
}

func mandatoryDetailFlags() {
	if common.ConnectorID == "" {
		connectorIDFlagPrompt()
	}
}

func optionalInitFlags(cmd *cobra.Command, args []string) {
	if ipfsAPIPort == 5001 {
		ipfsGatewayPortFlagPrompt()
	}
	if ipfsGatewayPort == 8080 {
		ipfsAPIPortFlagPrompt()
	}
	generalPrompt(cmd, args, promptStepSummary)
}

func connectorNameFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Connector Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	connectorName = result
}

func connectorTypeFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Connector Type",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	connectorType = result
}

func applicationIDFlagPrompt() {
	common.RequireApplication()
}

func networkIDFlagPrompt() {
	common.RequirePublicNetwork()
}

func organizationIDFlagPrompt() {
	common.RequireOrganization()
}

func connectorIDFlagPrompt() {
	common.RequireConnector(map[string]interface{}{})
}

func ipfsAPIPortFlagPrompt() {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "IPFS API Port",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	// turn to a uint somehow
	ipfsAPIPort, _ = strconv.ParseUint(result, 10, 10)
}

func ipfsGatewayPortFlagPrompt() {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "IPFS Gateway Port",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	ipfsGatewayPort, _ = strconv.ParseUint(result, 10, 10)
}
