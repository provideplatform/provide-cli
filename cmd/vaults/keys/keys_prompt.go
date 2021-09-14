package keys

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/provideplatform/provide-cli/cmd/shell"
	"github.com/spf13/cobra"

	vault "github.com/provideplatform/provide-go/api/vault"
)

var promptArgs []string

const promptStepInit = "Initialize"
const promptStepList = "List"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	if common.VaultID == "" {
		common.RequireVault()
	}

	switch step := currentStep; step {
	case promptStepInit:
		promptInit(cmd, args)
		createKeyRun(cmd, args)
	case promptStepList:
		promptList(cmd, args)
		listKeysRun(cmd, args)
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

func promptInit(cmd *cobra.Command, args []string) {
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.OrganizationID == "" {
		organizationidFlagPrompt()
	}
	if keytype == "" {
		keyTypePrompt()
	}
	if keyusage == "" {
		keyUsagePrompt()
	}
	if name == "" {
		nameFlagPrompt()
	}
	if description == "" {
		descriptionFlagPrompt()
	}
	if keyspec == "" {
		keySpecPrompt()
	}
}

func promptList(cmd *cobra.Command, args []string) {
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.OrganizationID == "" {
		organizationidFlagPrompt()
	}
	common.PromptPagination(paginate, pagination)
}

func optionalFlagsList(cmd *cobra.Command, args []string) {
	fmt.Println("Optional Flags:")
	if common.ApplicationID == "" {
		applicationIDFlagPrompt()
	}
	if common.OrganizationID == "" {
		applicationIDFlagPrompt()
	}
}

// Optional Flags For Init Key
func nameFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	name = result
}

func descriptionFlagPrompt() {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Description",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return
	}

	description = result
}

func keySpecPrompt() {
	prompt := promptui.Select{
		Label: "Spec",
		Items: []string{
			vault.KeySpecAES256GCM,
			vault.KeySpecECCBabyJubJub,
			vault.KeySpecChaCha20,
			vault.KeySpecECCC25519,
			vault.KeySpecECCBIP39,
			vault.KeySpecECCEd25519,
			vault.KeySpecECCSecp256k1,
			vault.KeySpecRSA2048,
			vault.KeySpecRSA3072,
			vault.KeySpecRSA4096,
		},
	}

	shell.MarshalPromptIO(&prompt)
	_, result, err := prompt.Run()
	if err != nil {
		return
	}

	keyspec = result
}

func keyTypePrompt() {
	prompt := promptui.Select{
		Label: "Type",
		Items: []string{"symmetric", "asymmetric"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	keytype = result
}

func keyUsagePrompt() {
	prompt := promptui.Select{
		Label: "Usage",
		Items: []string{"encrypt/decrypt", "sign/verify"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	keyusage = result
}

// Optional Flag For List Keys
func applicationIDFlagPrompt() {
	common.RequireApplication()
}

func organizationidFlagPrompt() {
	common.RequireOrganization()
}
