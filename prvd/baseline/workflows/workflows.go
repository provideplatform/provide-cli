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

package workflows

import (
	"github.com/provideplatform/provide-cli/prvd/baseline/workflows/messages"
	"github.com/provideplatform/provide-cli/prvd/baseline/workflows/worksteps"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

var paginate bool

var WorkflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "Interact with a baseline workflows",
	Long:  `Create, manage and interact with workflows via the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

func init() {
	WorkflowsCmd.AddCommand(listBaselineWorkflowsCmd)
	WorkflowsCmd.AddCommand(detailBaselineWorkflowCmd)
	WorkflowsCmd.AddCommand(initBaselineWorkflowCmd)
	WorkflowsCmd.AddCommand(deployBaselineWorkflowCmd)
	WorkflowsCmd.AddCommand(versionBaselineWorkflowCmd)
	WorkflowsCmd.AddCommand(worksteps.WorkstepsCmd)
	WorkflowsCmd.AddCommand(messages.MessagesCmd)
	WorkflowsCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
