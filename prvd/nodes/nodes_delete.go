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

var nodesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific node",
	Long:  `Delete a specific node by identifier and teardown any associated infrastructure`,
	Run:   deleteNode,
}

func deleteNode(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepDelete)
}

func deleteNodeRun(cmd *cobra.Command, args []string) {
	// FIXME!!!

	// token := common.RequireAPIToken()
	// status, _, err := provide.DeleteNetworkNode(token, common.NetworkID, common.NodeID)
	// if err != nil {
	// 	log.Printf("Failed to delete node with id: %s; %s", common.NodeID, err.Error())
	// 	os.Exit(1)
	// }
	// if status != 204 {
	// 	log.Printf("Failed to delete node with id: %s; received status: %d", common.NodeID, status)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Deleted node with id: %s", common.NodeID)
}

func init() {
	nodesDeleteCmd.Flags().StringVar(&common.NetworkID, "network", "", "network id")
	nodesDeleteCmd.MarkFlagRequired("network")

	nodesDeleteCmd.Flags().StringVar(&common.NodeID, "node", "", "id of the node")
	nodesDeleteCmd.MarkFlagRequired("node")
}
