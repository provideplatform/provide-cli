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

package baseline

import (
	"github.com/spf13/cobra"

	"github.com/provideplatform/provide-cli/prvd/baseline/stack"
	"github.com/provideplatform/provide-cli/prvd/baseline/subject_accounts"
	"github.com/provideplatform/provide-cli/prvd/baseline/workflows"
	"github.com/provideplatform/provide-cli/prvd/baseline/workgroups"
	"github.com/provideplatform/provide-cli/prvd/common"
)

var Optional bool

var BaselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Interact with the baseline protocol",
	Long:  `Interact with the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "See `prvd baseline stack --help` instead",
	Long: `Create, manage and interact with local baseline stack instances.

See: prvd baseline stack --help instead. This command is deprecated and will be removed soon.`,
	Run: func(cmd *cobra.Command, args []string) {
		generalPrompt(cmd, args, "")
	},
}

func init() {
	BaselineCmd.AddCommand(proxyCmd)
	BaselineCmd.AddCommand(stack.StackCmd)
	BaselineCmd.AddCommand(workgroups.WorkgroupsCmd)
	BaselineCmd.AddCommand(workflows.WorkflowsCmd)
	BaselineCmd.AddCommand(subject_accounts.SubjectAccountsCmd)
	BaselineCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the optional flags")
}
