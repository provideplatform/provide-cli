package nodes

import (
	"github.com/provideservices/provide-cli/cmd/common"
	// provide "github.com/provideservices/provide-go/api/nchain"

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
