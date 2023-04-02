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

package users

import (
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

var ParticipantsUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Interact with axiom user participants",
	Long:  `Create, manage and interact with user participants via the axiom protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

func init() {
	ParticipantsUsersCmd.AddCommand(inviteBaselineWorkgroupUserCmd)
	ParticipantsUsersCmd.AddCommand(listBaselineWorkgroupUsersCmd)

	ParticipantsUsersCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
	ParticipantsUsersCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
