package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"
	"github.com/spf13/cobra"
)

var nodeName string
var nodesInitCmd = &cobra.Command{
	Use:   "init --name Node1 --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new node",
	Long:  `Initialize a new node with options`,
	Run:   createNode,
}

func createNode(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{
		"name":      nodeName,
		"NetworkID": networkID,
		// "application_id": applicationID,
		// "address":        "0x",
		// "params":         contractParamsFactory(),
	}
	status, resp, err := provide.CreateNetworkNode(token, networkID, params)
	if err != nil {
		log.Printf("Failed to initialize node; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		node = resp.(map[string]interface{})
		nodeID = node["id"].(string)
		result := fmt.Sprintf("%s\t%s\n", node["id"], node["name"])
		fmt.Print(result)
	}
}

func init() {
	nodesInitCmd.Flags().StringVar(&nodeName, "name", "", "name of the node")
	nodesInitCmd.MarkFlagRequired("name")

	nodesInitCmd.Flags().StringVar(&networkID, "network", "", "target network id")
	nodesInitCmd.MarkFlagRequired("network")
}
