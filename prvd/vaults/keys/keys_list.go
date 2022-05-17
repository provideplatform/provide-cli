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

package keys

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	vault "github.com/provideplatform/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var keysListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of keys",
	Long:  `Retrieve a list of keys scoped to the authorized API token`,
	Run:   listKeys,
}

func listKeys(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listKeysRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	if common.OrganizationID != "" {
		params["organization_id"] = common.OrganizationID
	}
	resp, err := vault.ListKeys(token, common.VaultID, params)
	if err != nil {
		log.Printf("failed to retrieve keys list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp {
		vlt := resp[i]
		result := fmt.Sprintf("%s\t%s\t%s\n", vlt.ID.String(), *vlt.Name, *vlt.Description)
		fmt.Print(result)
	}
}

func init() {
	keysListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter keys")
	keysListCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier to filter keys")
	keysListCmd.Flags().StringVar(&common.VaultID, "vault", "", "identifier of the vault")

	keysListCmd.Flags().StringVar(&keyspec, "spec", "", "key spec query; non-matching keys are filtered")
	keysListCmd.Flags().StringVar(&keytype, "type", "", "key type query; non-matching keys are filtered")
	keysListCmd.Flags().StringVar(&keyusage, "usage", "", "key usage query; non-matching keys are filtered")

	keysListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	keysListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	keysListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of keys to retrieve per page")
}
