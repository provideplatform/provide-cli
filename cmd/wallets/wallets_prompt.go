package wallets

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case "Empty":
		emptyPrompt(cmd, args)
	case "Initialize":
		custodyPrompt(cmd, args)
	case "Custody":
		if flagPrompt() {
			optionalFlagsInit()
		}
	case "List":
		if flagPrompt() {
			optionalFlagsList()
		}
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{"Initialize", "List"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

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

func optionalFlagsInit() {
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

func optionalFlagsList() {
	fmt.Println("Optional Flags:")
	if common.ApplicationID == "" {
		applicationidFlagPrompt()
	}
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == "Initialize" {
		createManagedWallet(cmd, args)
	}
	if promptArgs[0] == "List" {
		listPrompt(cmd, args)
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
		fmt.Printf("Prompt failed %v\n", err)
	}

	generalPrompt(cmd, args, "Custody")
}

// Optional Flags For Init Wallet
//TODO: This is not custody theres has to be a better name
func custodialFlagPrompt() {
	prompt := promptui.Select{
		Label: "Would you like your wallet to be non-custodial",
		Items: []string{"Custodial", "Non-custodial"},
	}

	_, result, err := prompt.Run()

	nonCustodial = result != "Custodial"

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
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
		fmt.Printf("Prompt failed %v\n", err)
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

	// purpose, _ = strconv.ParseInt(result, 0, 64)
	// TODO: get rid of this
	purpose, _ = strconv.Atoi(result)

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

// List Wallets
func listPrompt(cmd *cobra.Command, args []string) {
	listWallets(cmd, args)
}

// Optional Flag For List Wallet
func applicationidFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Application ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	common.ApplicationID = result
}
