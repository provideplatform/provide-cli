package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/provideservices/provide-cli/cmd/accounts"
	"github.com/provideservices/provide-cli/cmd/api_tokens"
	"github.com/provideservices/provide-cli/cmd/applications"
	"github.com/provideservices/provide-cli/cmd/baseline"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/provideservices/provide-cli/cmd/connectors"
	"github.com/provideservices/provide-cli/cmd/contracts"
	"github.com/provideservices/provide-cli/cmd/networks"
	"github.com/provideservices/provide-cli/cmd/nodes"
	"github.com/provideservices/provide-cli/cmd/organizations"
	"github.com/provideservices/provide-cli/cmd/shell"
	"github.com/provideservices/provide-cli/cmd/users"
	"github.com/provideservices/provide-cli/cmd/vaults"
	"github.com/provideservices/provide-cli/cmd/wallets"
)

var rootCmd = &cobra.Command{
	Use:   "prvd",
	Short: "Provide CLI",
	Long: fmt.Sprintf(`%s

The Provide CLI exposes low-code tools to manage network, application and organization resources.

Run with the --help flag to see available options`, common.ASCIIBanner),
}

// Execute the default command path
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(common.InitConfig)

	rootCmd.PersistentFlags().BoolVarP(&common.Verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&common.CfgFile, "config", "c", "", "config file (default is $HOME/.provide-cli.yaml)")

	rootCmd.AddCommand(accounts.AccountsCmd)
	rootCmd.AddCommand(api_tokens.APITokensCmd)
	rootCmd.AddCommand(applications.ApplicationsCmd)
	rootCmd.AddCommand(users.AuthenticateCmd)
	rootCmd.AddCommand(baseline.BaselineCmd)
	rootCmd.AddCommand(connectors.ConnectorsCmd)
	rootCmd.AddCommand(contracts.ContractsCmd)
	rootCmd.AddCommand(networks.NetworksCmd)
	rootCmd.AddCommand(nodes.NodesCmd)
	rootCmd.AddCommand(organizations.OrganizationsCmd)
	rootCmd.AddCommand(shell.ShellCmd)
	rootCmd.AddCommand(users.UsersCmd)
	rootCmd.AddCommand(vaults.VaultsCmd)
	rootCmd.AddCommand(wallets.WalletsCmd)

	common.CacheCommands(rootCmd)
}
