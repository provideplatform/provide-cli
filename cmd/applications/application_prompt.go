package applications

import (
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepDetails = "Details"
const promptStepInit = "Initialize"
const promptStepList = "List"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		mandatoryInitFlag()
		if flagPrompt() {
			optionalFlagsInit()
		}
	case promptStepDetails:
		if flagPrompt() {
			optionalFlagsDetails()
		}
	case promptStepList:
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepDetails, promptStepList, promptStepInit},
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

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == promptStepInit {
		createApplication(cmd, args)
	}
	if promptArgs[0] == promptStepList {
		listApplications(cmd, args)
	}
	if promptArgs[0] == promptStepDetails {
		fetchApplicationDetails(cmd, args)
	}
}

func mandatoryInitFlag() {
	if applicationName == "" {
		applicationNameFlagPrompt()
	}
	if common.NetworkID == "" {
		NetworkIDFlagPrompt()
	}
}

func optionalFlagsInit() {
	if applicationType == "" {
		applicationTypeFlagPrompt()
	}
	if !baseline {
		baselineFlagPrompt()
	}
	if !withoutAccount {
		withoutAccountFlagPrompt()
	}
	if !withoutWallet {
		withoutWalletFlagPrompt()
	}
}

func optionalFlagsDetails() {
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
}

func applicationNameFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	applicationName = result
}

// NetworkIDFlagPrompt -- should we just use the common.RequireNetwork() convention instead of wrapping like this?
func NetworkIDFlagPrompt() {
	common.RequireNetwork()
}

func applicationTypeFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application Type",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	applicationName = result
}

func baselineFlagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"No", "Yes"},
	}

	_, flagResult, err := flagPrompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	baseline = flagResult == "Yes"
}

func withoutWalletFlagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"No", "Yes"},
	}

	_, result, err := flagPrompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	withoutWallet = result == "Yes"
}

func withoutAccountFlagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"No", "Yes"},
	}

	_, result, err := flagPrompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	withoutAccount = result == "Yes"
}

func applicationIDFlagPrompt() {
	common.RequireApplication()
}
