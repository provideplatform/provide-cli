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

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		mandatoryInitFlags()
		if flagPrompt() {
			optionalInitFlags()
		}
	case promptStepList:
		if flagPrompt() {
			applicationIDFlagPrompt()
		}
	case promptStepDetails:
		mandatoryDetailFlags()
	case promptStepDelete:
		mandatoryDeleteFlags()
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
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

func flagPrompt() bool {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"No", "Yes"},
	}

	_, flagResult, err := flagPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return false
	}

	return flagResult == "Yes"
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

func optionalInitFlags() {
	if ipfsAPIPort == 5001 {
		connectorNameFlagPrompt()
	}
	if ipfsGatewayPort == 8080 {
		connectorTypeFlagPrompt()
	}
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
		fmt.Printf("Prompt Exit\n")
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
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	connectorType = result
}

func applicationIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.ApplicationID = result
}

func networkIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Network ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.NetworkID = result
}

func connectorIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Connector ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.ConnectorID = result
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
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	// turn to a uint somehow
	ipfsAPIPort = result
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
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	// turn to a uint somehow
	ipfsGatewayPort = result
}
