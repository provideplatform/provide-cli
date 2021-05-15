package vaults

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepInit = "Initialize"
const promptStepList = "List"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		if flagPrompt() {
			optionalFlagsInit()
		}
	case promptStepList:
		if flagPrompt() {
			optionalFlagsList()
		}
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepInit, promptStepList},
	}

	_, result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

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
		os.Exit(1)
		return false
	}

	return flagResult == "Yes"
}

func optionalFlagsInit() {
	fmt.Println("Optional Flags:")
	if description == "" {
		descriptionFlagPrompt()
	}
	if name == "" {
		nameFlagPrompt()
	}
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.OrganizationID == "" {
		organizationidFlagPrompt()
	}

}

func optionalFlagsList() {
	fmt.Println("Optional Flags:")
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.OrganizationID == "" {
		applicationIDFlagPrompt()
	}
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == promptStepInit {
		createVault(cmd, args)
	}
	if promptArgs[0] == promptStepList {
		listVaults(cmd, args)
	}
}

// Optional Flags For Init Vault
func nameFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Vault Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	name = result
}

func descriptionFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Vault Description",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	description = result
}

// Optional Flag For List Vaults
func applicationIDFlagPrompt() {
	common.RequireApplication()
}

func organizationidFlagPrompt() {
	common.RequireOrganization()
}
