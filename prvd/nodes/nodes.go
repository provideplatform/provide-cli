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

package nodes

import (
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

var node map[string]interface{}
var nodes []interface{}
var nodeType string

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage nodes",
	Long:  `Manage and provision elastic distributed nodes`,
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
	NodesCmd.AddCommand(nodesInitCmd)
	NodesCmd.AddCommand(nodesLogsCmd)
	NodesCmd.AddCommand(nodesDeleteCmd)
	NodesCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")

}
