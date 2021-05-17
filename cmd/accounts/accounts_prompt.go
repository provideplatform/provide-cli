package accounts

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepCustody = "Custody"
const promptStepInit = "Initialize"
const promptStepList = "List"
const promptStepSummary = "Summary"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		custodyPrompt(cmd, args)
	case promptStepCustody:
		if flagPrompt(cmd, args) {
			optionalFlagsCustody(cmd, args)
		}
	case promptStepList:
		if flagPrompt(cmd, args) {
			optionalFlagsList(cmd, args)
		}
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

func optionalFlagsCustody(cmd *cobra.Command, args []string) {
	fmt.Println("Optional Flags:")
	if !nonCustodial {
		custodialFlagPrompt()
	}
	if accountName == "" {
		nameFlagPrompt()
	}
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.OrganizationID == "" {
		organizationIDFlagPrompt()
	}
	generalPrompt(cmd, args, promptStepSummary)
}

func optionalFlagsList(cmd *cobra.Command, args []string) {
	fmt.Println("Optional Flags:")
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	generalPrompt(cmd, args, promptStepSummary)
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == promptStepInit {
		CreateAccount(cmd, args)
	}
	if promptArgs[0] == promptStepList {
		listAccounts(cmd, args)
	}
}

// Init Wallet
func custodyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What type of Wallet would you like to create",
		Items: []string{"Managed", "Decentralised"},
	}

	_, result, err := prompt.Run()

	promptArgs = append(promptArgs, result)

	if err != nil {
		os.Exit(1)
		return
	}

	generalPrompt(cmd, args, promptStepCustody)
}

// Optional Flags For Init Wallet
//TODO: This is not custody theres has to be a better name
func custodialFlagPrompt() {
	prompt := promptui.Select{
		Label: "Would you like your wallet to be non-custodial",
		Items: []string{"Custodial", "Non-custodial"},
	}

	_, result, err := prompt.Run()

	nonCustodial = result != "Custodial"

	if err != nil {
		os.Exit(1)
		return
	}
}

func nameFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Wallet Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	accountName = result
}

func applicationIDFlagPrompt() {
	common.RequireApplication()
}

func networkIDFlagPrompt() {
	common.RequireNetwork()
}

func organizationIDFlagPrompt() {
	common.RequireOrganization()
}
