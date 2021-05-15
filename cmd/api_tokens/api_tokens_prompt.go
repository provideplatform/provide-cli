package api_tokens

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
		Items: []string{promptStepList, promptStepInit},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
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
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return false
	}

	return flagResult == "Yes"
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == promptStepInit {
		createAPIToken(cmd, args)
	}
	if promptArgs[0] == promptStepList {
		listAPITokens(cmd, args)
	}
}

func optionalFlagsInit() {
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.OrganizationID == "" {
		organizationIDFlagPrompt()
	}
	if !offlineAccess {
		offlineAccessFlagPrompt()
	}
	if !refreshToken {
		refreshTokenFlagPrompt()
	}
}

func optionalFlagsList() {
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
}

func refreshTokenFlagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"No", "Yes"},
	}

	_, result, err := flagPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	refreshToken = result == "Yes"
}

func offlineAccessFlagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"No", "Yes"},
	}

	_, result, err := flagPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	offlineAccess = result == "Yes"
}

func applicationIDFlagPrompt() {
	common.RequireApplication()
}

func organizationIDFlagPrompt() {
	common.RequireOrganization()
}
