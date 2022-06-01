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

package api_tokens

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var apiTokensListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of API tokens",
	Long:  `Retrieve a list of API tokens scoped to the authorized API token`,
	Run:   listAPITokens,
}

func listAPITokens(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	resp, err := provide.ListTokens(token, params)
	if err != nil {
		log.Printf("Failed to retrieve API tokens list; %s", err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve API tokens list; received status: %d", status)
	// 	os.Exit(1)
	// }
	for i := range resp {
		apiToken := resp[i]
		result := fmt.Sprintf("%s\t%s\n", apiToken.ID.String(), *apiToken.Token)
		fmt.Print(result)
	}
}

func init() {
	apiTokensListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter API tokens")
	apiTokensListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	apiTokensListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	apiTokensListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	apiTokensListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of API tokens to retrieve per page")
}
