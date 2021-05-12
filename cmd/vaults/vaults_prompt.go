package vaults

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case "Empty":
		emptyPrompt(cmd, args)
	case "Initialize":
		if flagPrompt() {
			optionalFlagsInit()
		}
	case "List":
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
		Items: []string{"Initialize", "List"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
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
		fmt.Printf("Prompt failed %v\n", err)
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
		applicationidFlagPrompt()
	}
	if common.OrganizationID == "" {
		organizationidFlagPrompt()
	}

}

func optionalFlagsList() {
	fmt.Println("Optional Flags:")
	if common.ApplicationID == "" {
		applicationidFlagPrompt()
	}
	if common.OrganizationID == "" {
		applicationidFlagPrompt()
	}
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == "Initialize" {
		createVault(cmd, args)
	}
	if promptArgs[0] == "List" {
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
		fmt.Printf("Prompt failed %v\n", err)
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
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	description = result
}

// Optional Flag For List Vaults
func applicationidFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	common.ApplicationID = result
}

func organizationidFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Organization ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	common.OrganizationID = result
}