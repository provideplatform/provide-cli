package api_tokens

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepInit = "Initialize"
const promptStepList = "List"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		if flagPrompt() {
			optionalFlagsInit()
		}
		createAPIToken(cmd, args)
	case promptStepList:
		if flagPrompt() {
			optionalFlagsList()
		}
		listAPITokens(cmd, args)
	case "":
		emptyPrompt(cmd, args)
	}
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepList, promptStepInit},
	}

	_, result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

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
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.OrganizationID == "" {
		organizationIDFlagPrompt()
	}
	if !refreshToken {
		refreshTokenFlagPrompt()
	}
	if !offlineAccess {
		offlineAccessFlagPrompt()
	}
	if refreshToken && offlineAccess {
		fmt.Println("⚠️  WARNING: You currently have both refresh and offline token set, Refresh token will take precedence")
	}
}

func optionalFlagsList() {
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
}

func refreshTokenFlagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to refresh your access token?",
		Items: []string{"No", "Yes"},
	}

	_, result, err := flagPrompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	refreshToken = result == "Yes"
}

func offlineAccessFlagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to vend an access/refresh token pair?",
		Items: []string{"No", "Yes"},
	}

	_, result, err := flagPrompt.Run()

	if err != nil {
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
