package messages

import (
	"errors"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepSend = "Send"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepSend:
		prompt()
		sendMessageRun(cmd, args)
	case "":
		emptyPrompt(cmd, args)
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

func prompt() {
	if common.ApplicationID == "" {
		common.RequireWorkgroup()
	}
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if messageType == "" {
		messageTypeFlagPrompt()
	}
	if id == "" {
		idFlagPrompt()
	}
	if baselineID == "" {
		baselineIDFlagPrompt()
	}
	if data == "" {
		dataFlagPrompt()
	}
}

// Optional Flag
func baselineIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Baseline ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	baselineID = result
}

// Mandatory Flags
func dataFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
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

func idFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	id = result
}

func messageTypeFlagPrompt() {
	if messageType != "" {
		return
	}

	items := map[string]string{
		"General Consistency": "general_consistency",
	}

	opts := make([]string, 0)
	for k := range items {
		opts = append(opts, k)
	}

	prmpt := promptui.Select{
		Label: "Message Type",
		Items: opts,
	}

	_, result, err := prmpt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	messageType = items[result]
}
