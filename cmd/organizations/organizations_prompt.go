package organizations

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/provideservices/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepDetails = "Details"
const promptStepInit = "Initialize"
const promptStepList = "List"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		nameFlagPrompt()
	case promptStepList:
		listOrganizations(cmd, args)
		return // FIXME
	case promptStepDetails:
		organizationFlagPrompt()
	default:
		emptyPrompt(cmd, args)
		return // FIXME
	}

	summary(cmd, args, promptArgs) // FIXME?
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepInit, promptStepList, promptStepDetails},
	}

	_, result, err := prompt.Run()

	promptArgs = append(promptArgs, result)

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	generalPrompt(cmd, args, result)
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == promptStepInit {
		createOrganization(cmd, args)
	}
	if promptArgs[0] == promptStepList {
		listOrganizations(cmd, args)
	}
	if promptArgs[0] == promptStepDetails {
		fetchOrganizationDetails(cmd, args)
	}
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
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	organizationName = result
}

// flag, equivalent to --organization
func organizationFlagPrompt() {
	opts := make([]string, 0)
	orgs, _ := ident.ListOrganizations(common.RequireUserAuthToken(), map[string]interface{}{})
	for _, org := range orgs {
		opts = append(opts, *org.Name)
	}

	prompt := promptui.Select{
		Label: "Select an organization:",
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		panic(err) // this will be recovered and we will exit...
	}

	fmt.Printf("selected %s at index: %v", *orgs[i].Name, i)
	common.OrganizationID = orgs[i].ID.String()
}
