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

package networks

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var networksDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a specific network",
	Long:  `Disable a specific network by identifier`,
	Run:   disableNetwork,
}

func disableNetwork(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	err := provide.UpdateNetwork(token, common.NetworkID, map[string]interface{}{
		"enabled": false,
	})
	if err != nil {
		log.Printf("Failed to disable network with id: %s; %s", common.NetworkID, err.Error())
		os.Exit(1)
	}
	// if status != 204 {
	// 	log.Printf("Failed to disable network with id: %s; received status: %d", common.NetworkID, status)
	// 	os.Exit(1)
	// }
	fmt.Printf("Disabled network with id: %s", common.NetworkID)
}

func init() {
	networksDisableCmd.Flags().StringVar(&common.NetworkID, "network", "", "id of the network")
	// networksDisableCmd.MarkFlagRequired("network")
}
