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

package accounts

import (
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

var AccountID string

var AccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage signing identities & accounts",
	Long: `Various APIs are exposed to provide convenient access to elliptic-curve cryptography
(ECC) helper methods such as generating managed (custodial) keypairs.

For convenience, it is also possible to generate keypairs with this utility which you (or your application)
is then responsible for securing. You should securely store any keys generated using this API. If you are
looking for hierarchical deterministic support, check out the wallets API.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")

		defer func() {
			if r := recover(); r != nil {
				os.Exit(1)
			}
		}()
	},
}

func init() {
	AccountsCmd.AddCommand(accountsListCmd)
	AccountsCmd.AddCommand(accountsInitCmd)
	AccountsCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	AccountsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
