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

package invitations

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var listBaselineWorkgroupInvitationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List workgroup invitations",
	Long:  `List the pending invitations for a baseline workgroup`,
	Run:   listInvitations,
}

func listInvitations(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listInvitationsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	common.AuthorizeOrganizationContext(false)

	token := common.RequireOrganizationToken()

	invitations, err := ident.ListApplicationInvitations(token, common.WorkgroupID, map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	})
	if err != nil {
		log.Printf("failed to fetch baseline workgroup invitations; %s", err.Error())
		os.Exit(1)
	}

	if len(invitations) == 0 {
		fmt.Print("No Pending Invitations Found")
	}

	for _, invitation := range invitations {
		result := fmt.Sprintf("%s\n", invitation.Email) // TODO-- make this show more relevant information
		fmt.Print(result)
	}
}

func init() {
	listBaselineWorkgroupInvitationsCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	listBaselineWorkgroupInvitationsCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")

	listBaselineWorkgroupInvitationsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineWorkgroupInvitationsCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of participants to retrieve per page")
	listBaselineWorkgroupInvitationsCmd.Flags().BoolVarP(&Optional, "Optional", "", false, "List all the Optional flags")
	listBaselineWorkgroupInvitationsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
