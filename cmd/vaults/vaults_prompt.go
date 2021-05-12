package vaults

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

// General Endpoints
func generalWalletPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case "empty":
		emptyWalletPrompt(cmd, args)
	case "init":
		flagPrompt()
	case "list":
		flagPrompt()
	default:
		emptyWalletPrompt(cmd, args)
	}
}

func emptyWalletPrompt(cmd *cobra.Command, args []string) {
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

	flagPrompt()

	summary(cmd, args, promptArgs)
}

func flagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"Set Optional Flags", "Dont Set Optional Flags"},
	}

	_, flagResult, err := flagPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if flagResult == "Set Optional Flags" {
		promptArgs = append(promptArgs, flagResult)
		if promptArgs[0] == "Initialize" {
			optionalFlagsInit()

		}
		if promptArgs[0] == "List" {
			optionalFlagsList()
		}
	}
}

func optionalFlagsInit() {
	fmt.Println("Optional Flags:")
	if description == "" {
		descriptionFlagVaultPrompt()
	}
	if name == "" {
		nameFlagVaultPrompt()
	}
	if common.ApplicationID == "" {
		applicationidFlagVaultPrompt()
	}
	if common.OrganizationID == "" {
		organizationidFlagPrompt()
	}

}

func optionalFlagsList() {
	fmt.Println("Optional Flags:")
	if common.ApplicationID == "" {
		applicationidFlagVaultPrompt()
	}
	if common.OrganizationID == "" {
		applicationidFlagVaultPrompt()
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

// Init Wallet
func custodyWalletPrompt(cmd *cobra.Command, args []string) string {
	prompt := promptui.Select{
		Label: "What type of Wallet would you like to create",
		Items: []string{"Managed", "Decentralised"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "nil"
	}

	return result
}

// Optional Flags For Init Wallet
func nameFlagVaultPrompt() {
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

func descriptionFlagVaultPrompt() {
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
func applicationidFlagVaultPrompt() {
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
