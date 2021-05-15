package contracts

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepExecute = "Execute"
const promptStepList = "List"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepExecute:
		mandatoryExecuteFlags()
		if flagPrompt() {
			optionalExecuteFlags()
		}
	case promptStepList:
		if flagPrompt() {
			optionalListFlags()
		}
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepExecute, promptStepList},
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
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return false
	}

	return flagResult == "Yes"
}
func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == promptStepExecute {
		executeContract(cmd, args)
	}
	if promptArgs[0] == promptStepList {
		listContracts(cmd, args)
	}
}

func mandatoryExecuteFlags() {
	if contractExecMethod == "" {
		methodFlagPrompt()
	}
	if common.ContractID == "" {
		contractIDFlagPrompt()
	}
}

func optionalExecuteFlags() {
	if contractExecMethod == "" {
		methodFlagPrompt()
	}
	if common.AccountID == "" {
		accountIDFlagPrompt()
	}
	if common.WalletID == "" {
		walletIDFlagPrompt()
	}
	if contractExecValue == 0 {
		valueFlagPrompt()
	}
}

func optionalListFlags() {
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
}

func methodFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	contractExecMethod = result
}

func contractIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.ContractID = result
}

func applicationIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.ApplicationID = result
}

func walletIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.WalletID = result
}

func accountIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.WalletID = result
}

func valueFlagPrompt() {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Wallet Purpose",
		Validate: validate,
	}

	result, err := prompt.Run()
	// Same issue as in networks
	contractExecValue, _ = strconv.ParseUint(result, 10, 64)

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}
}
