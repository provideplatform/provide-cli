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

package applications

import (
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const applicationTypeMessageBus = "message_bus"

var application map[string]interface{}

var ApplicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Manage applications",
	Long: `Create and manage logical applications which target a specific network and expose the following APIs:

	- API Tokens
	- Smart Contracts
	- Token Contracts
	- Signing Identities (wallets)
	- Oracles
	- Bridges
	- Connectors (i.e., IPFS)
	- Payment Hubs
	- Transactions`,
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
	ApplicationsCmd.AddCommand(applicationsListCmd)
	ApplicationsCmd.AddCommand(applicationsInitCmd)
	ApplicationsCmd.AddCommand(applicationsDetailsCmd)
	ApplicationsCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	ApplicationsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
