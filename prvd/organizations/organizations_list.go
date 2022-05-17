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

var page uint64
var rpp uint64

var organizationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of organizations",
	Long:  `Retrieve a list of organizations scoped to the authorized API token`,
	Run:   listOrganizations,
}

func listOrganizations(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listOrganizationsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	organizations, err := provide.ListOrganizations(token, params)
	if err != nil {
		log.Printf("Failed to retrieve organizations list; %s", err.Error())
		os.Exit(1)
	}
	for i := range organizations {
		organization := organizations[i]
		address := "0x"
		if addr, addrOk := organization.Metadata["address"].(string); addrOk {
			address = addr
		}
		result := fmt.Sprintf("%s\t%s\t%s\n", *organization.ID, *organization.Name, address)
		fmt.Print(result)
	}
}

func init() {
	organizationsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	organizationsListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	organizationsListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of organizations to retrieve per page")
}
