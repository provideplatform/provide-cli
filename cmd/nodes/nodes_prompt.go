package nodes

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var promptArgs []string

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case "Empty":
		emptyPrompt(cmd, args)
	case "Initialize":
		// mandatoryLogsFlags()
	case "Logs":
		//mandatoryLogsFlags()
		if flagPrompt() {
			//	optionalFlagsInit()
		}
	case "Delete":
		//mandatoryDeleteFlags()
		if flagPrompt() {
			//	optionalFlagsList()
		}
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{"Initialize", "Logs", "Delete"},
	}

	_, result, _ := prompt.Run()

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
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}

	return flagResult == "Yes"
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == "Initialize" {
		CreateNode(cmd, args)
	}
	if promptArgs[0] == "Logs" {
		nodeLogs(cmd, args)
	}
	if promptArgs[0] == "Delete" {
		nodeLogs(cmd, args)
	}
}
