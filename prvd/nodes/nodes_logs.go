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

package nodes

import (
	"github.com/provideplatform/provide-cli/prvd/common"
	// provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var nodesLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve logs for a node",
	Long:  `Retrieve paginated log output for a specific node by identifier`,
	Run:   nodeLogs,
}

func nodeLogs(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepLogs)
}

func nodeLogsRun(cmd *cobra.Command, args []string) {
	// FIXME
	// token := common.RequireAPIToken()
	// resp, err := provide.GetNetworkNodeLogs(token, common.NetworkID, common.NodeID, map[string]interface{}{
	// 	"page": page,
	// 	"rpp":  rpp,
	// })
	// if err != nil {
	// 	log.Printf("Failed to retrieve node logs for node with id: %s; %s", common.NodeID, err.Error())
	// 	os.Exit(1)
	// }
	// if status != 200 {
	// 	log.Printf("Failed to retrieve node logs for node with id: %s; received status: %d", common.NodeID, status)
	// 	os.Exit(1)
	// }
	// logsResponse := resp.(map[string]interface{})
	// if logs, logsOk := logsResponse["logs"].([]interface{}); logsOk {
	// 	for _, log := range logs {
	// 		fmt.Printf("%s\n", log)
	// 	}
	// }
	// if nextToken, nextTokenOk := logsResponse["next_token"].(string); nextTokenOk {
	// 	fmt.Printf("next token: %s", nextToken)
	// }
}

func init() {
	nodesLogsCmd.Flags().StringVar(&common.NetworkID, "network", "", "network id")
	nodesLogsCmd.MarkFlagRequired("network")

	nodesLogsCmd.Flags().StringVar(&common.NodeID, "node", "", "id of the node")
	nodesLogsCmd.MarkFlagRequired("node")

	nodesLogsCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	nodesLogsCmd.MarkFlagRequired("page")

	nodesLogsCmd.Flags().Uint64Var(&rpp, "rpp", 100, "number of log events to retrieve per page")
	nodesLogsCmd.MarkFlagRequired("rpp")
}
