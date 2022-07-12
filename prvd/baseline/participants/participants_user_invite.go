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
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	ident "github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var firstName string
var lastName string
var email string

var permissions int
var invitorAddress string
var registryContractAddress string
var managedTenant bool
var Optional bool
var paginate bool

var inviteBaselineWorkgroupUserCmd = &cobra.Command{
	Use:   "invite-user",
	Short: "Invite a user to a baseline workgroup",
	Long: `Invite a user to participate in a baseline workgroup.

A verifiable credential is issued which can then be distributed to the invited party out-of-band.`,
	Run: inviteUser,
}

func inviteUser(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInviteUser)
}

func inviteUserRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}
	if firstName == "" {
		firstNamePrompt()
	}
	if lastName == "" {
		lastNamePrompt()
	}
	if email == "" {
		emailPrompt()
	}

	common.AuthorizeOrganizationContext(false)

	token := common.RequireOrganizationToken()

	inviteParams := map[string]interface{}{
		"first_name":        firstName,
		"last_name":         lastName,
		"email":             email,
		"organization_name": common.Organization.Name,
		"application_id":    common.WorkgroupID, // FIXME--
		"params": map[string]interface{}{
			"workgroup":                   common.Workgroup,
			"is_organization_user_invite": true,
		},
	}

	if err := ident.CreateInvitation(token, inviteParams); err != nil {
		log.Printf("failed to invite baseline workgroup user; %s", err.Error())
		os.Exit(1)
	}

	log.Printf("invited baseline workgroup user: %s\n", email)
}

func firstNamePrompt() {
	prompt := promptui.Prompt{
		Label: "Invitee First Name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("first name required")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	firstName = result
}

func lastNamePrompt() {
	prompt := promptui.Prompt{
		Label: "Invitee Last Name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("last name required")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	lastName = result
}

func emailPrompt() {
	prompt := promptui.Prompt{
		Label: "Invitee Email",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("email required")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	email = result
}

func init() {
	inviteBaselineWorkgroupUserCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	inviteBaselineWorkgroupUserCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	inviteBaselineWorkgroupUserCmd.Flags().StringVar(&common.SubjectAccountID, "subject-account", "", "subject account identifier")
	inviteBaselineWorkgroupUserCmd.Flags().StringVar(&firstName, "first-name", "", "first name of the invited participant")
	inviteBaselineWorkgroupUserCmd.Flags().StringVar(&lastName, "last-name", "", "last name of the invited participant")
	inviteBaselineWorkgroupUserCmd.Flags().StringVar(&email, "email", "", "email address of the invited participant")

	// inviteBaselineWorkgroupUserCmd.Flags().BoolVar(&managedTenant, "managed-tenant", false, "if set, the invited participant is authorized to leverage operator-provided infrastructure")
	// inviteBaselineWorkgroupUserCmd.Flags().IntVar(&permissions, "permissions", 0, "permissions for invited participant")
	inviteBaselineWorkgroupUserCmd.Flags().BoolVarP(&Optional, "Optional", "", false, "List all the Optional flags")
}
