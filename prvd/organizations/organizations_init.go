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
	"github.com/provideplatform/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var organizationName string
var paginate bool

var organizationsInitCmd = &cobra.Command{
	Use:   "init --name 'Acme Inc.'",
	Short: "Initialize a new organization",
	Long:  `Initialize a new organization`,
	Run:   createOrganization,
}

func organizationConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"network_id": common.NetworkID,
	}

	return cfg
}
func createOrganization(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func createOrganizationRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"name":   organizationName,
		"config": organizationConfigFactory(),
	}
	organization, err := ident.CreateOrganization(token, params)
	if err != nil {
		log.Printf("Failed to initialize organization; %s", err.Error())
		os.Exit(1)
	}

	common.OrganizationID = *organization.ID

	orgToken, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to initialize organization; %s", err.Error())
		os.Exit(1)
	}

	if _, err := vault.CreateVault(*orgToken.AccessToken, map[string]interface{}{
		"name": fmt.Sprintf("%s vault", organizationName),
	}); err != nil {
		log.Printf("failed to create organization vault; %s", err.Error())
		os.Exit(1)
	}

	log.Printf("initialized organization: %s\t%s\n", organizationName, common.OrganizationID)
}

func init() {
	organizationsInitCmd.Flags().StringVar(&organizationName, "name", "", "name of the organization")
}
