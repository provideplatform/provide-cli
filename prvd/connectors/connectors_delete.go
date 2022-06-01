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

package connectors

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var connectorsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific connector",
	Long:  `Delete a specific connector by identifier and teardown any associated infrastructure`,
	Run:   deleteConnector,
}

func deleteConnector(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	err := provide.DeleteConnector(token, common.ConnectorID)
	if err != nil {
		log.Printf("Failed to delete connector with id: %s; %s", common.ConnectorID, err.Error())
		os.Exit(1)
	}
	// if status != 204 {
	// 	log.Printf("Failed to delete connector with id: %s; received status: %d", common.ConnectorID, status)
	// 	os.Exit(1)
	// }
	fmt.Printf("Deleted connector with id: %s", common.ConnectorID)
}

func init() {
	connectorsDeleteCmd.Flags().StringVar(&common.ConnectorID, "connector", "", "id of the connector")
	// connectorsDeleteCmd.MarkFlagRequired("connector")

	connectorsDeleteCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application id")
	// connectorsDeleteCmd.MarkFlagRequired("application")
}
