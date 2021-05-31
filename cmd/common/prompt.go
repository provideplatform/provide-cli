package common

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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

func normaliseCmd(cmd *cobra.Command, args []string) (string, string) {
	flag, _ := regexp.Compile("\\[(.*)")
	r, _ := regexp.Compile("\\--(.*)")
	usedCommand := strings.Split(cmd.UseLine(), flag.FindString(cmd.UseLine()))
	normalisedCommand := strings.TrimSpace(strings.Join(usedCommand, ""))
	argsLine := strings.TrimSpace(strings.Join(args, " "))

	rmFlagsLine := strings.Split(normalisedCommand, r.FindString(normalisedCommand))

	normalisedOutput := strings.TrimSpace(strings.Join(rmFlagsLine, ""))

	return normalisedOutput, argsLine
}

func CacheCommands(cmd *cobra.Command) {
	if commands == nil {
		commands = map[string]*cobra.Command{}
	}

	command, _ := normaliseCmd(cmd, nil)
	hashCmd := sha256.Sum256([]byte(command))

	commands[string(hashCmd[:])] = cmd
	for _, child := range cmd.Commands() {
		CacheCommands(child)
	}
}

func CmdExists(cmd *cobra.Command, args []string) (bool, string) {
	command, arguments := normaliseCmd(cmd, args)

	argsCommandNormalised := fmt.Sprintf("%s%s", command, arguments)

	hashCmd := sha256.Sum256([]byte(argsCommandNormalised))

	return commands[string(hashCmd[:])] != nil, argsCommandNormalised
}

func CmdExistsOrExit(cmd *cobra.Command, args []string) {
	exists, command := CmdExists(cmd, args)
	if !exists {
		fmt.Println(command + " is not a valid command")
		os.Exit(1)
	}
}

// RequireApplication is equivalent to a required --application flag
func RequireApplication() error {
	if ApplicationID != "" {
		return nil
	}

	opts := make([]string, 0)
	apps, _ := ident.ListApplications(RequireUserAccessToken(), map[string]interface{}{})
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
	apps, _ := ident.ListApplications(RequireUserAccessToken(), map[string]interface{}{
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
	orgs, _ := ident.ListOrganizations(RequireUserAccessToken(), map[string]interface{}{})
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
		Label: requireAccountSelectLabel,
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

var MandatoryValidation = func(input string) error {
	if len(input) < 1 {
		return errors.New("password must have more than 6 characters")
	}
	return nil
}

var MandatoryNumberValidation = func(input string) error {
	if len(input) < 1 {
		return errors.New("field cant be nil")
	}
	_, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}

var NumberValidation = func(input string) error {
	_, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}

var NoValidation = func(input string) error {
	return nil
}

var JSONValidation = func(input string) error {
	if len(input) < 1 {
		return errors.New("field cant be nil")
	}

	var js map[string]interface{}
	if json.Unmarshal([]byte(input), &js) != nil {
		return errors.New("invalid JSON")

	}
	return nil
}

var HexValidation = func(input string) error {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	isHex := re.MatchString(input)

	if !isHex {
		return errors.New("input is not a Hex")
	}
	return nil
}

func FreeInput(label string, defaultValue string, validate func(string) error) string {

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
