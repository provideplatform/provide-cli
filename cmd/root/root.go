package root

import (
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

func Root(rootCmd *cobra.Command, enableShell bool) {
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
	if enableShell {
		rootCmd.AddCommand(shell.ShellCmd)
	}
	rootCmd.AddCommand(users.UsersCmd)
	rootCmd.AddCommand(vaults.VaultsCmd)
	rootCmd.AddCommand(wallets.WalletsCmd)
}
