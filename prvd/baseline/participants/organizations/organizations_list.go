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

package organizations

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

var listBaselineWorkgroupOrganizationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List workgroup organizations",
	Long:  `List the organizations for a baseline workgroup`,
	Run:   listOrganizations,
}

func listOrganizations(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listOrganizationsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	common.AuthorizeOrganizationContext(false)

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to fetch baseline workgroup organizations; %s", err.Error())
		os.Exit(1)
	}

	orgs, err := ident.ListApplicationOrganizations(*token.AccessToken, common.WorkgroupID, map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	})
	if err != nil {
		log.Printf("failed to fetch baseline workgroup organizations; %s", err.Error())
		os.Exit(1)
	}

	for _, org := range orgs {
		result := fmt.Sprintf("%s\t%s\n", *org.ID, *org.Name) // TODO-- show DegreeOfSeparation
		fmt.Print(result)
	}
}

func init() {
	listBaselineWorkgroupOrganizationsCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	listBaselineWorkgroupOrganizationsCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")

	listBaselineWorkgroupOrganizationsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineWorkgroupOrganizationsCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of participants to retrieve per page")
	listBaselineWorkgroupOrganizationsCmd.Flags().BoolVarP(&Optional, "Optional", "", false, "List all the Optional flags")
	listBaselineWorkgroupOrganizationsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
