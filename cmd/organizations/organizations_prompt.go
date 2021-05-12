package organizations

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

// General Endpoints
func generalOrganizationPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case "Empty":
		emptyOrganizationPrompt(cmd, args)
	case "Initialize":
		promptArgs = append(promptArgs, "Initialize")
		mandatoryFlagsInit()
	case "List":
		promptArgs = append(promptArgs, "List")
		summary(cmd, args, promptArgs)
	case "Details":
		promptArgs = append(promptArgs, "Details")
		mandatoryFlagsDetails()
	default:
		emptyOrganizationPrompt(cmd, args)
	}
}

func emptyOrganizationPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{"Initialize", "List", "Details"},
	}

	_, result, err := prompt.Run()

	generalOrganizationPrompt(cmd, args, result)

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	promptArgs = append(promptArgs, result)

	summary(cmd, args, promptArgs)
}

func mandatoryFlagsInit() {
	nameFlagOrganizationPrompt()
}

func mandatoryFlagsDetails() {
	organizationidFlagPrompt()
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == "Initialize" {
		createOrganization(cmd, args)
	}
	if promptArgs[0] == "List" {
		listOrganizations(cmd, args)
	}
	if promptArgs[0] == "Details" {
		fetchOrganizationDetails(cmd, args)
	}
}

// Mandatory Flags For Init Wallet
func nameFlagOrganizationPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Organization Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	organizationName = result
}

// Mandatory Flag For detail Organizations
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
