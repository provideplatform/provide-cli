package wallets

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

const promptStepCustody = "Custody"
const promptStepInit = "Initialize"
const promptStepList = "List"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		custodyPrompt(cmd, args)
	case promptStepCustody:
		if flagPrompt(cmd, args) {
			optionalFlagsInit(cmd, args)
		}
		CreateWalletRun(cmd, args)
	case promptStepList:
		if flagPrompt(cmd, args) {
			optionalFlagsList(cmd, args)
		}
		listWalletsRun(cmd, args)
	case "":
		emptyPrompt(cmd, args)
	}
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepInit, promptStepList},
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
	return flagResult == "Yes"
}

func optionalFlagsInit(cmd *cobra.Command, args []string) {
	fmt.Println("Optional Flags:")
	if !nonCustodial {
		custodialFlagPrompt()
	}
	if walletName == "" {
		nameFlagPrompt()
	}
	if purpose == 44 {
		purposeFlagPrompt()
	}
}

func optionalFlagsList(cmd *cobra.Command, args []string) {
	fmt.Println("Optional Flags:")
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
}

// Init Wallet
func custodyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What type of Wallet would you like to create",
		Items: []string{"Managed", "Decentralised"},
	}

	_, result, err := prompt.Run()

	promptArgs = append(promptArgs, result)

	if err != nil {
		os.Exit(1)
		return
	}

	generalPrompt(cmd, args, promptStepCustody)
}

// Optional Flags For Init Wallet
func custodialFlagPrompt() {
	prompt := promptui.Select{
		Label: "Would you like your wallet to be non-custodial",
		Items: []string{"Custodial", "Non-custodial"},
	}

	_, result, err := prompt.Run()

	nonCustodial = result != "Custodial"

	if err != nil {
		os.Exit(1)
		return
	}
}

func nameFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Wallet Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	walletName = result
}

func purposeFlagPrompt() {
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

	purpose, _ = strconv.Atoi(result)

	if err != nil {
		os.Exit(1)
		return
	}
}

// Optional Flag For List Wallet
func applicationIDFlagPrompt() {
	common.RequireApplication()
}
