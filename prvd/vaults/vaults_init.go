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

package vaults

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var name string
var description string
var optional bool
var paginate bool

var vaultsInitCmd = &cobra.Command{
	Use:   "init --name 'My Vault' --description 'not your keys, not your crypto'",
	Short: "Create a new vault",
	Long:  `Initialize a new vault`,
	Run:   createVault,
}

func createVault(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func createVaultRun(cmd *cobra.Command, args []string) {

	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"name":        name,
		"description": description,
	}
	vlt, err := provide.CreateVault(token, params)
	if err != nil {
		log.Printf("Failed to genereate HD wallet; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\t%s\n", vlt.ID.String(), *vlt.Name, *vlt.Description)
	fmt.Print(result)
}

func init() {
	vaultsInitCmd.Flags().StringVar(&name, "name", "", "name of the vault")
	vaultsInitCmd.Flags().StringVar(&description, "description", "", "description of the vault")

	vaultsInitCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier for which the vault will be created")
	vaultsInitCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier for which the vault will be created")
	vaultsInitCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
