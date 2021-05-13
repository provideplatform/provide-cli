package nodes

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

const promptStepLogs = "Logs"
const promptStepInit = "Initialize"
const promptStepDelete = "Delete"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepInit:
		mandatoryInitFlags()
		if flagPrompt() {
			optionalFlagsInit()
		}
	case promptStepDelete:
		mandatoryDeleteFlags()
	case promptStepLogs:
		mandatoryLogFlags()
	default:
		emptyPrompt(cmd, args)
	}

	summary(cmd, args, promptArgs)
}

func emptyPrompt(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "What would you like to do",
		Items: []string{promptStepInit, promptStepLogs, promptStepDelete},
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
	if promptArgs[0] == promptStepInit {
		CreateNode(cmd, args)
	}
	if promptArgs[0] == promptStepLogs {
		nodeLogs(cmd, args)
	}
	if promptArgs[0] == promptStepDelete {
		deleteNode(cmd, args)
	}
}

func mandatoryInitFlags() {
	if common.NetworkID == "" {
		networkIDFlagPrompt()
	}
	if common.Image == "" {
		imageFlagPrompt()
	}
	if role == "" {
		roleFlagPrompt()
	}
}

func mandatoryLogFlags() {
	if common.NetworkID == "" {
		networkIDFlagPrompt()
	}
	if common.NodeID == "" {
		nodeIDFlagPrompt()
	}
	if page == 1 {
		pageFlagPrompt()
	}
	if rpp == 100 {
		roleFlagPrompt()
	}
}

func mandatoryDeleteFlags() {
	if common.NetworkID == "" {
		networkIDFlagPrompt()
	}
	if common.NodeID == "" {
		nodeIDFlagPrompt()
	}
}

func optionalFlagsInit() {
	fmt.Println("Optional Flags:")
	if common.HealthCheckPath == "" {
		healthCheckPathFlagPrompt()
	}
	if common.TCPIngressPorts == "" {
		TCPIngressPortsFlagPrompt()
	}
	if common.UDPIngressPorts == "" {
		UDPIngressPortsFlagPrompt()
	}
	if common.TaskRole == "" {
		taskRoleFlagPrompt()
	}
}

// Flags
func networkIDFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Network ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.NetworkID = result
}

func imageFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Image",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.Image = result
}

func roleFlagPrompt() {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("name cant be nil")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Role",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	role = result
}

func healthCheckPathFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Health Check Path",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.HealthCheckPath = result
}

func TCPIngressPortsFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "TCP Ingress Ports",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.TCPIngressPorts = result
}

func UDPIngressPortsFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "UDP Ingress Ports",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.UDPIngressPorts = result
}

func taskRoleFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Task Role",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.UDPIngressPorts = result
}

func nodeIDFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Node ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}

	common.NodeID = result
}

func pageFlagPrompt() {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Node ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}
	page, _ = strconv.Atoi(result)
}

func rppFlagPrompt() {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Node ID",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt Exit\n")
		os.Exit(1)
		return
	}
	rpp, _ = strconv.ParseUint(result, 10, 64)

}
