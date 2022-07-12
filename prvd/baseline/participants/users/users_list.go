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
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var listBaselineWorkgroupUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List workgroup users",
	Long:  `List the users for a baseline workgroup`, // TODO-- actually lists organization users
	Run:   listUsers,
}

func listUsers(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listUsersRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}

	common.AuthorizeOrganizationContext(false)

	token := common.RequireOrganizationToken()

	users, err := ident.ListOrganizationUsers(token, common.OrganizationID, map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	})
	if err != nil {
		log.Printf("failed to fetch baseline workgroup users; %s", err.Error())
		os.Exit(1)
	}

	for _, user := range users {
		result := fmt.Sprintf("%s\t%s\n", *user.ID, user.Name) // TODO-- show role from permissions / Workgroup.UserID
		fmt.Print(result)
	}
}

func init() {
	listBaselineWorkgroupUsersCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")

	listBaselineWorkgroupUsersCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineWorkgroupUsersCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of participants to retrieve per page")
	listBaselineWorkgroupUsersCmd.Flags().BoolVarP(&Optional, "Optional", "", false, "List all the Optional flags")
	listBaselineWorkgroupUsersCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
