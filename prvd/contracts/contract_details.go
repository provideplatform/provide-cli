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

package contracts

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var contractsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific smart contract",
	Long:  `Retrieve details for a specific smart contract by identifier, scoped to the authorized API token`,
	Run:   fetchContractDetails,
}

func fetchContractDetails(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	contract, err := provide.GetContractDetails(token, common.ContractID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for contract with id: %s; %s", common.ContractID, err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve details for contract with id: %s; %s", common.ContractID, resp)
	// 	os.Exit(1)
	// }
	result := fmt.Sprintf("%s\t%s\n", contract.ID.String(), *contract.Name)
	fmt.Print(result)
}

func init() {
	contractsDetailsCmd.Flags().StringVar(&common.ContractID, "contract", "", "id of the contract")
	contractsDetailsCmd.MarkFlagRequired("contract")
}
