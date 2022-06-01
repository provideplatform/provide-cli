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

package wallets

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

var walletsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of custodial HD wallets",
	Long:  `Retrieve a list of HD wallets scoped to the authorized API token`,
	Run:   listWallets,
}

func listWallets(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listWalletsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	resp, err := provide.ListWallets(token, params)
	if err != nil {
		log.Printf("Failed to retrieve wallets list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp {
		wallet := resp[i]
		result := fmt.Sprintf("%s\t%s\n", wallet.ID.String(), *wallet.PublicKey)
		// FIXME-- when wallet.Name exists... result = fmt.Sprintf("Wallet %s\t%s - %s\n", wallet.Name, wallet.ID.String(), *wallet.Address)
		fmt.Print(result)
	}
}

func init() {
	walletsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter HD wallets")
	walletsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	walletsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	walletsListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	walletsListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of wallets to retrieve per page")
}
