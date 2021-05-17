package users

import (
	"errors"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func firstNamePrompt() {
	prompt := promptui.Prompt{
		Label: "First Name",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	firstName = result
}

func lastNamePrompt() {
	prompt := promptui.Prompt{
		Label: "Last Name",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	lastName = result
}

func emailPrompt() {
	prompt := promptui.Prompt{
		Label: "Email",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	email = result
}

func passwordPrompt() {
	validate := func(input string) error {
		if len(input) < 8 {
			return errors.New("password must have at least 8 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     '*',
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	passwd = result
}

func usersPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What command would you like to run",
		Items: []string{"Authenticate", "Create"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	switch result {
	case "Authenticate":
		authenticate(cmd, args)
	case "Create":
		create(cmd, args)
	}
}
