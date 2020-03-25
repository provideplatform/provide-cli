package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var page uint
var rpp uint

var nodesLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve logs for a node",
	Long:  `Retrieve paginated log output for a specific node by identifier`,
	Run:   nodeLogs,
}

func nodeLogs(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	status, resp, err := provide.GetNetworkNodeLogs(token, networkID, nodeID, map[string]interface{}{
		"page": page,
		"rpp":  rpp,
	})
	if err != nil {
		log.Printf("Failed to retrieve node logs for node with id: %s; %s", nodeID, err.Error())
		os.Exit(1)
	}
	if status != 200 {
		log.Printf("Failed to retrieve node logs for node with id: %s; received status: %d", nodeID, status)
		os.Exit(1)
	}
	logsResponse := resp.(map[string]interface{})
	if logs, logsOk := logsResponse["logs"].([]interface{}); logsOk {
		for _, log := range logs {
			fmt.Printf("%s\n", log)
		}
	}
	if nextToken, nextTokenOk := logsResponse["next_token"].(string); nextTokenOk {
		fmt.Printf("next token: %s", nextToken)
	}
}

func init() {
	nodesLogsCmd.Flags().StringVar(&networkID, "network", "", "network id")
	nodesLogsCmd.MarkFlagRequired("network")

	nodesLogsCmd.Flags().StringVar(&nodeID, "node", "", "id of the node")
	nodesLogsCmd.MarkFlagRequired("node")

	nodesLogsCmd.Flags().UintVar(&page, "page", 1, "page number to retrieve")
	nodesLogsCmd.MarkFlagRequired("page")

	nodesLogsCmd.Flags().UintVar(&rpp, "rpp", 100, "number of log events to retrieve per page")
	nodesLogsCmd.MarkFlagRequired("rpp")
}
