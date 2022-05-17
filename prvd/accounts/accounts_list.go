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

package accounts

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"
	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of signing identities",
	Long:  `Retrieve a list of signing identities (accounts) scoped to the authorized API token`,
	Run:   listAccounts,
}

func listAccounts(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	resp, err := provide.ListAccounts(token, params)
	if err != nil {
		log.Printf("Failed to retrieve accounts list; %s", err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve accounts list; received status: %d", status)
	// 	os.Exit(1)
	// }
	for i := range resp {
		account := resp[i]
		result := fmt.Sprintf("%s\t%s\n", account.ID.String(), account.Address)
		// TODO-- when account.Name exists... result = fmt.Sprintf("%s\t%s - %s\n", name, account, *account.Address)
		fmt.Print(result)
	}
}

func init() {
	accountsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter accounts")
	accountsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	accountsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	accountsListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	accountsListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of accounts to retrieve per page")
}
