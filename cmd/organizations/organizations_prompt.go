package organizations

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepDetails = "Details"
const promptStepInit = "Initialize"
const promptStepList = "List"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, step string) {
	switch step {
	case promptStepInit:
		nameFlagPrompt()
		createOrganizationRun(cmd, args)
	case promptStepList:
		listOrganizationsRun(cmd, args)
	case promptStepDetails:
		organizationIDFlagPrompt()
		fetchOrganizationDetailsRun(cmd, args)
	case "":
		emptyPrompt(cmd, args)
	default:
		fmt.Println("no-ops")
	}
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do?",
		Items: []string{promptStepInit, promptStepList, promptStepDetails},
	}

	_, result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	generalPrompt(cmd, args, result)
}

// Mandatory Flags For Init Wallet
func nameFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Organization Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	organizationName = result
}

// require organization...
func organizationIDFlagPrompt() {
	common.RequireOrganization()
}
