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

package applications

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/accounts"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-cli/prvd/wallets"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var applicationName string
var applicationType string
var baseline bool
var withoutAPIToken bool
var withoutAccount bool
var withoutWallet bool
var optional bool
var paginate bool

var applicationsInitCmd = &cobra.Command{
	Use:   "init --name 'my app' --network 024ff1ef-7369-4dee-969c-1918c6edb5d4 [--baseline]",
	Short: "Initialize a new application",
	Long:  `Initialize a new application targeting a specified mainnet`,
	Run:   createApplication,
}

func applicationConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"network_id": common.NetworkID,
	}

	if applicationType != "" {
		cfg["type"] = applicationType
	}

	return cfg
}

func createApplication(cmd *cobra.Command, args []string) {
	if withoutAPIToken && !withoutWallet {
		fmt.Println("Cannot create an application that has a wallet but no API token.")
		os.Exit(1)
	}
	token := common.RequireAPIToken()
	cfg := applicationConfigFactory()
	if baseline {
		cfg["baseline"] = true
	}
	params := map[string]interface{}{
		"name":   applicationName,
		"type":   applicationType,
		"config": cfg,
	}

	application, err := provide.CreateApplication(token, params)
	if err != nil {
		log.Printf("Failed to initialize application; %s", err.Error())
		os.Exit(1)
	}

	// // FIXME-- authorize app token...
	// token := application.Token
	// common.ApplicationID = application.ID.String().(string)
	// applicationToken := token.Token

	// appAPITokenKey := common.BuildConfigKeyWithID(common.APITokenConfigKeyPartial, common.ApplicationID)
	// if !viper.IsSet(appAPITokenKey) {
	// 	viper.Set(appAPITokenKey, applicationToken)
	// 	viper.WriteConfig()
	// }
	// fmt.Printf("Application API Token\t%s\n", applicationToken)

	result := fmt.Sprintf("%s\t%s\n", application.ID.String(), *application.Name)
	fmt.Print(result)
	if !withoutAccount {
		accounts.CreateAccount(cmd, args)
	}
	if !withoutWallet {
		wallets.CreateWallet(cmd, args)
	}
}

func init() {
	applicationsInitCmd.Flags().StringVar(&applicationName, "name", "", "name of the application")
	// applicationsInitCmd.MarkFlagRequired("name")

	applicationsInitCmd.Flags().StringVar(&applicationType, "type", "", "application type")

	applicationsInitCmd.Flags().StringVar(&common.NetworkID, "network", "", "target network id")
	// applicationsInitCmd.MarkFlagRequired("network")

	applicationsInitCmd.Flags().BoolVar(&baseline, "baseline", false, "setup a baseline workgroup")

	applicationsInitCmd.Flags().BoolVar(&withoutAccount, "without-account", false, "do not create a new account (signing identity) for this application")
	applicationsInitCmd.Flags().BoolVar(&withoutWallet, "without-wallet", false, "do not create a new HD wallet for this application")
	applicationsInitCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")

}
