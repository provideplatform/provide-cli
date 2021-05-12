package users

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func emailPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Email",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	email = result
}

func passwordPrompt() {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("Password must have more than 6 characters")
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
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	passwd = result
}

func authenticatePrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What command would you like to run",
		Items: []string{"Authenticate"},
	}

	_, result, err := prompt.Run()

	if result == "Authenticate" {
		authenticate(cmd, args)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
}
