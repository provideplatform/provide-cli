package common

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/api/nchain"
)

const requireAccountSelectLabel = "Select an account:"
const requireApplicationSelectLabel = "Select an application:"
const requireConnectorSelectLabel = "Select a connector:"
const requireNetworkSelectLabel = "Select a network:"
const requireOrganizationSelectLabel = "Select an organization:"
const requireWalletSelectLabel = "Select a wallet:"
const requireWorkgroupSelectLabel = "Select a workgroup:"

// RequireApplication is equivalent to a required --application flag
func RequireApplication() error {
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

	fmt.Printf("selected application %s at index: %v", *apps[i].Name, i)
	ApplicationID = apps[i].ID.String()
	return nil
}

// RequireWorkgroup is equivalent to a required --workgroup flag
// (yes, this is identical to RequireApplication() with exception to the Printf content...)
func RequireWorkgroup() error {
	opts := make([]string, 0)
	apps, _ := ident.ListApplications(RequireUserAuthToken(), map[string]interface{}{})
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

	fmt.Printf("selected workgroup %s at index: %v", *apps[i].Name, i)
	ApplicationID = apps[i].ID.String()
	return nil
}

// RequireConnector is equivalent to a required --connector flag
func RequireConnector(params map[string]interface{}) error {
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

	fmt.Printf("selected connector %s at index: %v", *connectors[i].Name, i)
	ConnectorID = connectors[i].ID.String()
	return nil
}

// RequireNetwork is equivalent to a required --network flag
func RequireNetwork() error {
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

	fmt.Printf("selected network %s at index: %v", *networks[i].Name, i)
	NetworkID = networks[i].ID.String()
	return nil
}

// RequireOrganization is equivalent to a required --organization flag
func RequireOrganization() error {
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

	fmt.Printf("selected organization %s at index: %v", *orgs[i].Name, i)
	OrganizationID = orgs[i].ID.String()
	return nil
}

// RequireAccount is equivalent to a required --account flag
func RequireAccount(params map[string]interface{}) error {
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

	fmt.Printf("selected account %s at index: %v", *accounts[i].PublicKey, i)
	AccountID = accounts[i].ID.String()
	return nil
}

// RequireWallet is equivalent to a required --wallet flag
func RequireWallet() error {
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

	fmt.Printf("selected wallet %s at index: %v", *wallets[i].PublicKey, i)
	WalletID = wallets[i].ID.String()
	return nil
}
