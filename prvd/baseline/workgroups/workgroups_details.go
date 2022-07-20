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

package workgroups

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"

	"github.com/spf13/cobra"
)

var detailBaselineWorkgroupCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific baseline workgroup",
	Long:  `Retrieve details for a specific baseline workgroup by identifier, scoped to the authorized API token`,
	Run:   fetchWorkgroupDetails,
}

func fetchWorkgroupDetails(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepDetails)
}

func fetchWorkgroupDetailsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}

	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("Failed to retrieve details for workgroup with id: %s; %s", common.WorkgroupID, err.Error())
		os.Exit(1)
	}

	wg, err := baseline.GetWorkgroupDetails(*token.AccessToken, common.WorkgroupID, map[string]interface{}{})
	if err != nil {
		log.Printf("Failed to retrieve details for workgroup with id: %s; %s", common.WorkgroupID, err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(wg, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func init() {
	detailBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")
	detailBaselineWorkgroupCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
}
