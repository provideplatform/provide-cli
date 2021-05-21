package prompt

import (
	"os"

	"github.com/manifoldco/promptui"
)

func freeInput(lable string) string {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Baseline Id",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return err.Error()
	}

	return result
}

func selectInput(args []string, label string) string {
	prompt := promptui.Select{
		Label: label,
		Items: args,
	}

	_, result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return err.Error()
	}

	return result

}
