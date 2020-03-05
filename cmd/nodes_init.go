package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"
	"github.com/spf13/cobra"
)

var nodeName string
var imageID string
var roleID string
var taskRole string
var nodesInitCmd = &cobra.Command{
	Use:   "init --network 024ff1ef-7369-4dee-969c-1918c6edb5d4 --image redis --provider docker --region us-east-1 --role redis --target aws",
	Short: "Initialize a new node",
	Long:  `Initialize a new node with options`,
	Run:   createNode,
}

func nodeSecurityConfigFactory() map[string]interface{} {
	return map[string]interface{}{
		"health_check": map[string]interface{}{
			"path": "/api/v0/version",
		},
		"egress": "*",
		"ingress": map[string]interface{}{
			"0.0.0.0/0": map[string]interface{}{
				"tcp": []uint{ipfsAPIPort, ipfsGatewayPort},
				"udp": []uint{},
			},
		},
	}
}

func nodeConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"image":       imageID,
		"credentials": infrastructureCredentialsConfigFactory(),
		"p2p":         false,
		"region":      region,
		"target_id":   targetID,
		"task_role":   taskRole,
		"provider_id": providerID,
		"engine_id":   connectorType,
		"role":        roleID,
		"container":   container,
		"env":         map[string]interface{}{},
	}

	securityCfg := nodeSecurityConfigFactory()
	if securityCfg != nil {
		cfg["security"] = securityCfg
	}

	if connectorType == connectorTypeIPFS {
		if ipfsAPIPort != 0 {
			cfg["api_port"] = ipfsAPIPort
		}
		if ipfsGatewayPort != 0 {
			cfg["gateway_port"] = ipfsGatewayPort
		}
	}

	return cfg
}
func createNode(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{
		"network_id": networkID,
		"config":     nodeConfigFactory(),
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
	// nodesInitCmd.Flags().StringVar(&nodeName, "name", "", "name of the node")
	// nodesInitCmd.MarkFlagRequired("name")

	nodesInitCmd.Flags().StringVar(&networkID, "network", "", "target network id")
	nodesInitCmd.MarkFlagRequired("network")

	nodesInitCmd.Flags().StringVar(&imageID, "image", "", "image id")
	nodesInitCmd.MarkFlagRequired("image")

	nodesInitCmd.Flags().StringVar(&roleID, "role", "", "role id")
	nodesInitCmd.MarkFlagRequired("role")

	requireInfrastructureFlags(connectorsInitCmd, false)

	connectorsInitCmd.Flags().UintVar(&ipfsAPIPort, "ipfs-api-port", 5001, "tcp listen port for the ipfs api")
	connectorsInitCmd.Flags().UintVar(&ipfsGatewayPort, "ipfs-gateway-port", 8080, "tcp listen port for the ipfs gateway")
}
