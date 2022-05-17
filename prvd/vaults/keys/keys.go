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

package keys

import (
	"os"

	"github.com/spf13/cobra"
)

var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys",
	Long: `Create and manage cryptographic keys.

Supports symmetric and asymmetric key specs with encrypt/decrypt and sign/verify operations.

Docs: https://docs.provide.services/vault/api-reference/keys`,
	Run: func(cmd *cobra.Command, args []string) {
		generalPrompt(cmd, args, "")

		defer func() {
			if r := recover(); r != nil {
				os.Exit(1)
			}
		}()
	},
}

func init() {
	KeysCmd.AddCommand(keysListCmd)
	KeysCmd.AddCommand(keysInitCmd)
	KeysCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
