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

package subject_accounts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"

	"github.com/spf13/cobra"
)

var subjectAccountDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific subject account",
	Long:  `Retrieve details for a specific baseline subject account`,
	Run:   fetchSubjectAccountDetails,
}

func fetchSubjectAccountDetails(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepDetails)
}

// SHA256 is a convenience method to return the sha256 hash of the given input
func SHA256(str string) string {
	digest := sha256.New()
	digest.Write([]byte(str))
	return hex.EncodeToString(digest.Sum(nil))
}

func fetchSubjectAccountDetailsRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}

	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}

	token := common.RequireOrganizationToken()

	saID := SHA256(fmt.Sprintf("%s.%s", common.OrganizationID, common.WorkgroupID))

	sa, err := baseline.GetSubjectAccountDetails(token, common.OrganizationID, saID, map[string]interface{}{})
	if err != nil {
		log.Printf("Failed to retrieve details for subject account with id: %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}

	if sa.ID == nil {
		fmt.Print("subject account not found")
		return
	}

	result := fmt.Sprintf("%s;\tworkgroup: %s\t%s;\torganization: %s\t%s\n", *sa.ID, common.Workgroup.ID, *common.Workgroup.Name, *common.Organization.ID, *common.Organization.Name)
	fmt.Print(result)
}

func init() {
	subjectAccountDetailsCmd.Flags().StringVar(&common.SubjectAccountID, "subject account", "", "subject account identifier")
}
