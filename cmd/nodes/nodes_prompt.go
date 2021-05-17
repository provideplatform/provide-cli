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
		if flagPrompt(cmd, args) {
			optionalFlagsInit(cmd, args)
		}
		CreateNodeRun(cmd, args)
	case promptStepDelete:
		mandatoryDeleteFlags()
		deleteNodeRun(cmd, args)
	case promptStepLogs:
		mandatoryLogFlags()
		nodeLogsRun(cmd, args)
	case "":
		emptyPrompt(cmd, args)
	}
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

func flagPrompt(cmd *cobra.Command, args []string) bool {
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
		rppFlagPrompt()
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

func optionalFlagsInit(cmd *cobra.Command, args []string) {
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

func networkIDFlagPrompt() {
	common.RequirePublicNetwork()
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
		Label:    "Page",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}
	page, _ = strconv.ParseUint(result, 10, 64)
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
		Label:    "RPP",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}
	rpp, _ = strconv.ParseUint(result, 10, 64)
}
