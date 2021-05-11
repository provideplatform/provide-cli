package wallets

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

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
		fmt.Println("list")
	default:
		emptyWalletPrompt(cmd, args)
	}
}

func emptyWalletPrompt(cmd *cobra.Command, args []string) {
	var promptArgs []string
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
	if result == "List" {
		listWalletPrompt(cmd, args)
	}

	if result == "Initialize" {
		promptArgs = append(promptArgs, initWalletPrompt(cmd, args))
	}

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
		optionalFlagsInit()
	}

	summary(cmd, args, promptArgs)
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

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	fmt.Println(promptArgs[0])
	if promptArgs[0] == "Initialize" {
		fmt.Printf("Creating a %v wallet ", promptArgs[1])
		if !nonCustodial && walletName == "" && purpose == 44 {
			fmt.Printf("with no optional flags. Flags set to default values. \n")
		} else {
			fmt.Printf("with optional flags: \n")
			if nonCustodial {
				fmt.Printf("\tCustody: non-custodial\n")
			}
			if walletName != "" {
				fmt.Printf("\tWallet Name: %v \n", walletName)
			}
			if purpose != 44 {
				fmt.Printf("\tPurpose: %v \n", purpose)
			}
		}
		createManagedWallet(cmd, args)
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
func optionalFlagsList() {
}
