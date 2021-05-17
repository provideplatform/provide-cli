package messages

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepSend = "Send"
const promptStepSummary = "Summary"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepSend:
		mandatoryFlagsSend()
		if flagPrompt(cmd, args) {
			optionalFlagsSend(cmd, args)
		}
	case "":
		emptyPrompt(cmd, args)
	default:
		fmt.Println("no-ops")
	}
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepSend},
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

func optionalFlagsSend(cmd *cobra.Command, args []string) {

}

func mandatoryFlagsSend() {

}

// Mandatory Flags
func dataFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Data",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	data = result
}
