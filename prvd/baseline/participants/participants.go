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

package participants

import (
	participants_invitations "github.com/provideplatform/provide-cli/prvd/baseline/participants/invitations"
	participants_organizations "github.com/provideplatform/provide-cli/prvd/baseline/participants/organizations"
	participants_users "github.com/provideplatform/provide-cli/prvd/baseline/participants/users"
	"github.com/provideplatform/provide-cli/prvd/common"

	"github.com/spf13/cobra"
)

var ParticipantsCmd = &cobra.Command{
	Use:   "participants",
	Short: "Interact with participants in a baseline workgroup",
	Long:  `Invite, manage and interact with workgroup participants via the baseline protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdExistsOrExit(cmd, args)

		generalPrompt(cmd, args, "")
	},
}

func init() {
	ParticipantsCmd.AddCommand(participants_users.ParticipantsUsersCmd)
	ParticipantsCmd.AddCommand(participants_organizations.ParticipantsOrganizationsCmd)
	ParticipantsCmd.AddCommand(participants_invitations.ParticipantsInvitationsCmd)

	ParticipantsCmd.Flags().BoolVarP(&Optional, "Optional", "", false, "List all the Optional flags")
	ParticipantsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
