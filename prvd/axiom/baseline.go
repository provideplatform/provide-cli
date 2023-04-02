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

package axiom

import (
	"github.com/spf13/cobra"

	"github.com/provideplatform/provide-cli/prvd/axiom/domain_models"
	"github.com/provideplatform/provide-cli/prvd/axiom/stack"
	"github.com/provideplatform/provide-cli/prvd/axiom/subject_accounts"
	"github.com/provideplatform/provide-cli/prvd/axiom/workflows"
	"github.com/provideplatform/provide-cli/prvd/axiom/workgroups"
	"github.com/provideplatform/provide-cli/prvd/common"
)

var Optional bool

var BaselineCmd = &cobra.Command{
	Use:   "axiom",
	Short: "Interact with the axiom protocol",
	Long:  `Interact with the axiom protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "See `prvd axiom stack --help` instead",
	Long: `Create, manage and interact with local axiom stack instances.

See: prvd axiom stack --help instead. This command is deprecated and will be removed soon.`,
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
	BaselineCmd.AddCommand(domain_models.DomainModelsCmd)
	BaselineCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the optional flags")
}
