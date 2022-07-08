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
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var organizationsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific organization",
	Long:  `Retrieve details for a specific organization by identifier, scoped to the authorized API token`,
	Run:   fetchOrganizationDetails,
}

func fetchOrganizationDetails(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepDetails)
}

func fetchOrganizationDetailsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		generalPrompt(cmd, args, "Details")
	}
	token := common.RequireUserAccessToken()
	params := map[string]interface{}{}
	organization, err := provide.GetOrganizationDetails(token, common.OrganizationID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for organization with id: %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve details for organization with id: %s; %s", common.OrganizationID, organization)
	// 	os.Exit(1)
	// }
	result := fmt.Sprintf("%s\t%s\n", *organization.ID, *organization.Name)
	fmt.Print(result)
}

func init() {
	organizationsDetailsCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "id of the organization")
	// organizationsDetailsCmd.MarkFlagRequired("organization")
}
