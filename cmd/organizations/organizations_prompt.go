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
		summary(cmd, args, promptArgs)
	case promptStepDetails:
		organizationidFlagPrompt()
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
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

// Mandatory Flag For detail Organizations
func organizationidFlagPrompt() {
	orgs, _ := ident.ListOrganizations(common.RequireUserAuthToken(), map[string]interface{}{})
	// for _, org := range orgs {
	// 	orgNames = append(orgNames, *org.Name)
	// 	orgIDs = append(orgIDs, org.ID.String())
	// }

	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: orgs,
	}

	i, result, err := prompt.Run()

	fmt.Printf("Result: %v at index: %v", result, orgs[i].ID.String())
	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.OrganizationID = orgs[i].ID.String()
}
