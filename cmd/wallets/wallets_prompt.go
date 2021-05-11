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
func generalWalletPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case "empty":
		emptyWalletPrompt(cmd, args)
	case "init":
		initWalletPrompt(cmd, args)
	case "custodial flag":
		custodialFlagWalletPrompt()
	case "list":
		emptyWalletPrompt(cmd, args)
	default:
		emptyWalletPrompt(cmd, args)
	}
}

func emptyWalletPrompt(cmd *cobra.Command, args []string) {
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

	if result == "Initialize" {
		promptArgs = append(promptArgs, initWalletPrompt(cmd, args))
	}

	flagPrompt()

	summary(cmd, args, promptArgs)
}

func flagPrompt() {
	flagPrompt := promptui.Select{
		Label: "Would you like to set Optional Flags?",
		Items: []string{"Set Optional Flags", "Dont Set Optional Flags"},
	}

	_, flagResult, err := flagPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if flagResult == "Set Optional Flags" {
		promptArgs = append(promptArgs, flagResult)
		if promptArgs[0] == "Initialize" {
			optionalFlagsInit()

		}
		if promptArgs[0] == "List" {
			optionalFlagsList()

		}
	}
}

func optionalFlagsInit() {
	fmt.Println("Optional Flags:")
	if !nonCustodial {
		custodialFlagWalletPrompt()
	}
	if walletName == "" {
		nameFlagWalletPrompt()
	}
	if purpose == 44 {
		purposeFlagWalletPrompt()
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
		listWalletPrompt(cmd, args)
	}
}

// Init Wallet
func initWalletPrompt(cmd *cobra.Command, args []string) string {
	prompt := promptui.Select{
		Label: "What type of Wallet would you like to create",
		Items: []string{"Managed", "Decentralised"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "nil"
	}

	return result
}

// Optional Flags For Init Wallet
func custodialFlagWalletPrompt() {
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

func nameFlagWalletPrompt() {
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

func purposeFlagWalletPrompt() {
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
func listWalletPrompt(cmd *cobra.Command, args []string) {
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
