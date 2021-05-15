package networks

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

var promptArgs []string

const promptStepInit = "Initialize"
const promptStepList = "List"
const promptStepDisable = "Disable"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		mandatoryFlagsInit()
	case promptStepList:
		if flagPrompt() {
			publicFlagPrompt()
		}
	case promptStepDisable:
		networkIDFlagPrompt()
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepInit, promptStepList, promptStepDisable},
	}

	_, result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
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
		os.Exit(1)
		return false
	}

	return flagResult == "Yes"
}

func mandatoryFlagsInit() {
	if chain == "" {
		chainFlagPrompt()
	}
	if nativeCurrency == "" {
		nativeCurrencyFlagPrompt()
	}
	if platform == "" {
		platformFlagPrompt()
	}
	if protocolID == "" {
		protocolIDFlagPrompt()
	}
	if networkName == "" {
		networkNameFlagPrompt()
	}
}

func summary(cmd *cobra.Command, args []string, promptArgs []string) {
	if promptArgs[0] == promptStepInit {
		CreateNetwork(cmd, args)
	}
	if promptArgs[0] == promptStepList {
		listNetworks(cmd, args)
	}
	if promptArgs[0] == promptStepDisable {
		//listWallets(cmd, args)
	}
}

// Init Wallet
func chainFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Chain",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	chain = result
}

func nativeCurrencyFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Native Currency",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	nativeCurrency = result
}

func platformFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Platform",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	platform = result
}
func protocolIDFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Protocol ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	protocolID = result
}

func networkNameFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Network Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	networkName = result
}

func publicFlagPrompt() {
	prompt := promptui.Select{
		Label: "Would you like the network to be public",
		Items: []string{"No", "Yes"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	public = result == "Yes"
}

func networkIDFlagPrompt() {
	common.RequireNetwork()
}
