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

package contracts

import (
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const contractTypeRegistry = "registry"

var contract map[string]interface{}
var contracts []interface{}
var contractType string

var ContractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Manage smart contracts",
	Long:  `Compile and deploy smart contracts locally from source or execute previously-deployed contracts`,
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
	ContractsCmd.AddCommand(contractsListCmd)
	ContractsCmd.AddCommand(contractsExecuteCmd)
	ContractsCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	ContractsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
