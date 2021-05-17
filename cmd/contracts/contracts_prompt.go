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
const promptStepSummary = "Summary"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepExecute:
		mandatoryExecuteFlags()
		if flagPrompt(cmd, args) {
			optionalExecuteFlags(cmd, args)
		}
		return
	case promptStepList:
		if flagPrompt(cmd, args) {
			optionalListFlags(cmd, args)
		}
	case promptStepSummary:
		summary(cmd, args, promptArgs)
	case "":
		emptyPrompt(cmd, args)
	default:
		fmt.Println("no-ops")
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

func optionalExecuteFlags(cmd *cobra.Command, args []string) {
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
	summary(cmd, args, promptArgs)
}

func optionalListFlags(cmd *cobra.Command, args []string) {
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	summary(cmd, args, promptArgs)
}

func methodFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Method",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
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
		Label:    "Contract ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	common.ContractID = result
}

func accountIDFlagPrompt() {
	common.RequireAccount(map[string]interface{}{})
}

func applicationIDFlagPrompt() {
	common.RequireApplication()
}

func networkIDFlagPrompt() {
	common.RequireNetwork()
}

func walletIDFlagPrompt() {
	common.RequireWallet()
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
		os.Exit(1)
		return
	}
}
