/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package prvd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/provideplatform/provide-cli/prvd/accounts"
	"github.com/provideplatform/provide-cli/prvd/api_tokens"
	"github.com/provideplatform/provide-cli/prvd/applications"
	axiom "github.com/provideplatform/provide-cli/prvd/axiom"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-cli/prvd/connectors"
	"github.com/provideplatform/provide-cli/prvd/contracts"
	"github.com/provideplatform/provide-cli/prvd/networks"
	"github.com/provideplatform/provide-cli/prvd/nodes"
	"github.com/provideplatform/provide-cli/prvd/organizations"
	"github.com/provideplatform/provide-cli/prvd/shell"
	"github.com/provideplatform/provide-cli/prvd/users"
	"github.com/provideplatform/provide-cli/prvd/vaults"
	"github.com/provideplatform/provide-cli/prvd/wallets"
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
	rootCmd.AddCommand(axiom.BaselineCmd)
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
