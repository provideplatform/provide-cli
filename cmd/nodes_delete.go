package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var nodesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific node",
	Long:  `Delete a specific node by identifier and teardown any associated infrastructure`,
	Run:   deleteNode,
}

func deleteNode(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	status, _, err := provide.DeleteNetworkNode(token, networkID, nodeID)
	if err != nil {
		log.Printf("Failed to delete node with id: %s; %s", nodeID, err.Error())
		os.Exit(1)
	}
	if status != 204 {
		log.Printf("Failed to delete node with id: %s; received status: %d", nodeID, status)
		os.Exit(1)
	}
	fmt.Printf("Deleted node with id: %s", nodeID)
}

func init() {
	nodesDeleteCmd.Flags().StringVar(&networkID, "network", "", "network id")
	nodesDeleteCmd.MarkFlagRequired("network")

	nodesDeleteCmd.Flags().StringVar(&nodeID, "node", "", "id of the node")
	nodesDeleteCmd.MarkFlagRequired("node")
}
