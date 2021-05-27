package common

import (
	"errors"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/api/nchain"
	"github.com/provideservices/provide-go/api/vault"
	"github.com/spf13/cobra"
)

const requireAccountSelectLabel = "Select an account"
const requireApplicationSelectLabel = "Select an application"
const requireConnectorSelectLabel = "Select a connector"
const requireNetworkSelectLabel = "Select a network"
const requireOrganizationSelectLabel = "Select an organization"
const requireVaultSelectLabel = "Select a vault"
const requireWalletSelectLabel = "Select a wallet"
const requireWorkgroupSelectLabel = "Select a workgroup"

var commands map[string]*cobra.Command

// func CacheCommands(cmd *cobra.Command) {
// 	if commands == nil {
// 		commands = map[string]*cobra.Command{}
// 	}

// 	hashCmd := sha256.Sum256([]byte(cmd.UseLine()))

// 	commands[string(hashCmd[:])] = cmd
// 	for _, child := range cmd.Commands() {
// 		fmt.Print(child)
// 		CacheCommands(child)
// 	}
// }

// func CmdExists(cmd *cobra.Command) bool {
// 	hashCmd := sha256.Sum256([]byte(cmd.UseLine()))
// 	fmt.Println(hashCmd)

// 	return commands[string(hashCmd[:])] != nil

// }

// func CmdExistsOrExit(cmd *cobra.Command) {
// 	if !CmdExists(cmd) {
// 		os.Exit(1)
// 	}
// }

// RequireApplication is equivalent to a required --application flag
func RequireApplication() error {
	if ApplicationID != "" {
		return nil
	}

	opts := make([]string, 0)
	apps, _ := ident.ListApplications(RequireUserAuthToken(), map[string]interface{}{})
	for _, app := range apps {
		opts = append(opts, *app.Name)
	}

	prompt := promptui.Select{
		Label: requireApplicationSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	ApplicationID = apps[i].ID.String()
	return nil
}

// RequireWorkgroup is equivalent to a required --workgroup flag
// (yes, this is identical to RequireApplication() with exception to the Printf content...)
func RequireWorkgroup() error {
	if ApplicationID != "" {
		return nil
	}

	opts := make([]string, 0)
	apps, _ := ident.ListApplications(RequireUserAuthToken(), map[string]interface{}{
		"type": "baseline",
	})
	for _, app := range apps {
		opts = append(opts, *app.Name)
	}

	prompt := promptui.Select{
		Label: requireWorkgroupSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	ApplicationID = apps[i].ID.String()
	return nil
}

// RequireConnector is equivalent to a required --connector flag
func RequireConnector(params map[string]interface{}) error {
	if ConnectorID != "" {
		return nil
	}

	opts := make([]string, 0)
	connectors, _ := nchain.ListConnectors(RequireAPIToken(), params)
	for _, connector := range connectors {
		opts = append(opts, *connector.Name)
	}

	prompt := promptui.Select{
		Label: requireConnectorSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	ConnectorID = connectors[i].ID.String()
	return nil
}

// RequireNetwork is equivalent to a required --network flag
func RequireNetwork() error {
	if NetworkID != "" {
		return nil
	}

	opts := make([]string, 0)
	networks, _ := nchain.ListNetworks(RequireAPIToken(), map[string]interface{}{})
	for _, network := range networks {
		opts = append(opts, *network.Name)
	}

	prompt := promptui.Select{
		Label: requireNetworkSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	NetworkID = networks[i].ID.String()
	return nil
}

// RequirePublicNetwork is equivalent to a required --network flag; but list options filtered to show only public networks
func RequirePublicNetwork() error {
	if NetworkID != "" {
		return nil
	}

	opts := make([]string, 0)
	networks, _ := nchain.ListNetworks(RequireAPIToken(), map[string]interface{}{
		"public": "true",
	})
	for _, network := range networks {
		opts = append(opts, *network.Name)
	}

	prompt := promptui.Select{
		Label: requireNetworkSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	NetworkID = networks[i].ID.String()
	return nil
}

// RequireOrganization is equivalent to a required --organization flag
func RequireOrganization() error {
	if OrganizationID != "" {
		return nil
	}

	opts := make([]string, 0)
	orgs, _ := ident.ListOrganizations(RequireUserAuthToken(), map[string]interface{}{})
	for _, org := range orgs {
		opts = append(opts, *org.Name)
	}

	prompt := promptui.Select{
		Label: requireOrganizationSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	Organization = orgs[i]
	OrganizationID = orgs[i].ID.String()
	return nil
}

// RequireVault is equivalent to a required --vault flag
func RequireVault() error {
	if VaultID != "" {
		return nil
	}

	opts := make([]string, 0)
	vaults, _ := vault.ListVaults(RequireAPIToken(), map[string]interface{}{})
	for _, vlt := range vaults {
		opts = append(opts, *vlt.Name)
	}

	prompt := promptui.Select{
		Label: requireVaultSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	VaultID = vaults[i].ID.String()
	return nil
}

// RequireAccount is equivalent to a required --account flag
func RequireAccount(params map[string]interface{}) error {
	if AccountID != "" {
		return nil
	}

	opts := make([]string, 0)
	accounts, _ := nchain.ListAccounts(RequireAPIToken(), params)
	for _, acct := range accounts {
		opts = append(opts, *acct.PublicKey)
	}

	prompt := promptui.Select{
		Label: requireWalletSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	AccountID = accounts[i].ID.String()
	return nil
}

// RequireWallet is equivalent to a required --wallet flag
func RequireWallet() error {
	if WalletID != "" {
		return nil
	}

	opts := make([]string, 0)
	wallets, _ := nchain.ListWallets(RequireAPIToken(), map[string]interface{}{})
	for _, wallet := range wallets {
		opts = append(opts, *wallet.PublicKey)
	}

	prompt := promptui.Select{
		Label: requireWalletSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	WalletID = wallets[i].ID.String()
	return nil
}

func FreeInput(label string, defaultValue string, validation string) string {
	validate := func(input string) error {
		return nil
	}

	if validation == "Mandatory" {
		validate = func(input string) error {
			if len(input) < 1 {
				return errors.New("password must have more than 6 characters")
			}
			return nil
		}
	}

	var prompt = promptui.Prompt{}
	if label == "Password" {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: validate,
			Default:  defaultValue,
			Mask:     '*',
		}

	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: validate,
			Default:  defaultValue,
		}
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return err.Error()
	}

	return result
}

func SelectInput(args []string, label string) string {
	prompt := promptui.Select{
		Label: label,
		Items: args,
	}

	_, result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
		return err.Error()
	}

	return result

}
