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
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/axiom"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var listBaselineSubjectAccountsCmd = &cobra.Command{
	Use:   "list",
	Short: "List axiom subject accounts",
	Long:  `List all available axiom subject accounts`,
	Run:   listSubjectAccounts,
}

func listSubjectAccounts(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listSubjectAccountsRun(cmd *cobra.Command, args []string) {
	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to retrieve axiom subject accounts; %s", err.Error())
		os.Exit(1)
	}

	subject_accounts, err := axiom.ListSubjectAccounts(*token.AccessToken, common.OrganizationID, map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	})
	if err != nil {
		log.Printf("failed to retrieve axiom subject accounts; %s", err.Error())
		os.Exit(1)
	}
	// fmt.Printf("subject accounts len: %v", len(subject_accounts))
	for _, subject_account := range subject_accounts {
		details, err := axiom.GetSubjectAccountDetails(*token.AccessToken, common.OrganizationID, *subject_account.ID, map[string]interface{}{})
		if err != nil {
			log.Printf("failed to retrieve axiom subject accounts; %s", err.Error())
			os.Exit(1)
		}

		subject_account_wg, err := axiom.GetWorkgroupDetails(*token.AccessToken, *details.Metadata.WorkgroupID, map[string]interface{}{})
		if err != nil {
			log.Printf("failed to retrieve axiom subject accounts; %s", err.Error())
			os.Exit(1)
		}

		subject_account_org, err := ident.GetOrganizationDetails(*token.AccessToken, *details.Metadata.OrganizationID, map[string]interface{}{})
		if err != nil {
			log.Printf("failed to retrieve axiom subject accounts; %s", err.Error())
			os.Exit(1)
		}

		result := fmt.Sprintf("%s;\tworkgroup: %s\t%s;\torganization: %s\t%s\n", *subject_account.ID, subject_account_wg.ID, *subject_account_wg.Name, *subject_account_org.ID, *subject_account_org.Name)
		fmt.Print(result)
	}
}

func init() {
	listBaselineSubjectAccountsCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier")

	listBaselineSubjectAccountsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineSubjectAccountsCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of axiom subject accounts to retrieve per page")
	listBaselineSubjectAccountsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
