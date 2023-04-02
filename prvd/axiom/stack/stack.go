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

package stack

import (
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

var StackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Interact with a local axiom stack",
	Long:  `Create, manage and interact with local axiom stack instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

var runBaselineStackCmd = &cobra.Command{
	Use:   "run",
	Short: "See `prvd axiom stack start --help` instead",
	Long: `Start a local axiom stack instance and connect to internal systems of record.

See: prvd axiom stack run --help instead. This command is deprecated and will be removed soon.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStackStart(cmd, args)
	},
}

func init() {
	StackCmd.AddCommand(logsBaselineStackCmd)
	StackCmd.AddCommand(runBaselineStackCmd)
	StackCmd.AddCommand(startBaselineStackCmd)
	StackCmd.AddCommand(stopBaselineStackCmd)
	StackCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the optional flags")
}
