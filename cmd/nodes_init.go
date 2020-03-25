package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/provideservices/provide-go"
	"github.com/spf13/cobra"
)

var isP2P bool
var role string

var nodesInitCmd = &cobra.Command{
	Use:   "init --network 024ff1ef-7369-4dee-969c-1918c6edb5d4 --image redis --provider docker --region us-east-1 --role redis --target aws",
	Short: "Initialize a new node",
	Long:  `Initialize a new node with options`,
	Run:   CreateNode,
}

func nodeEnvConfigFactory() map[string]interface{} {
	return map[string]interface{}{}
}

func nodeSecurityConfigFactory() map[string]interface{} {
	tcpIngress := make([]uint, 0)
	udpIngress := make([]uint, 0)

	for _, port := range strings.Split(tcpIngressPorts, ",") {
		if port != "" {
			portInt, err := strconv.Atoi(port)
			if err != nil {
				log.Printf("Invalid tcp ingress port: %s", port)
				os.Exit(1)
			}
			tcpIngress = append(tcpIngress, uint(portInt))
		}
	}

	for _, port := range strings.Split(udpIngressPorts, ",") {
		if port != "" {
			portInt, err := strconv.Atoi(port)
			if err != nil {
				log.Printf("Invalid udp ingress port: %s", port)
				os.Exit(1)
			}
			udpIngress = append(udpIngress, uint(portInt))
		}
	}

	cfg := map[string]interface{}{
		"egress": "*",
		"ingress": map[string]interface{}{
			"0.0.0.0/0": map[string]interface{}{
				"tcp": tcpIngress,
				"udp": udpIngress,
			},
		},
	}

	var healthCheck map[string]interface{}
	if healthCheckPath != "" {
		healthCheck = map[string]interface{}{
			"path": healthCheckPath,
		}
	}
	if healthCheck != nil && len(healthCheck) > 0 {
		cfg["health_check"] = healthCheck
	}

	return cfg
}

func nodeConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"credentials": infrastructureCredentialsConfigFactory(),
		"entrypoint":  nil,
		"env":         nodeEnvConfigFactory(),
		"image":       image,
		"p2p":         isP2P,
		"region":      region,
		"role":        role,
		"target_id":   targetID,
	}

	if taskRole != "" {
		cfg["task_role"] = taskRole
	}

	// if resources != nil {
	// 	cfg["resources"] = resources
	// }

	securityCfg := nodeSecurityConfigFactory()
	if securityCfg != nil {
		cfg["security"] = securityCfg
	}

	return cfg
}

// CreateNode deploys a node to an existing peer-to-peer network;
// see https://docs.provide.services/microservices/goldmine/#deploy-network-node
func CreateNode(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{
		"config": nodeConfigFactory(),
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
	nodesInitCmd.Flags().StringVar(&networkID, "network", "", "target network id")
	nodesInitCmd.MarkFlagRequired("network")

	nodesInitCmd.Flags().BoolVar(&isP2P, "p2p", true, "when true, genesis state and peer resolution are enforced during initialization")

	nodesInitCmd.Flags().StringVar(&image, "image", "", "docker image; can be an official image name or fully-qualified repo")
	nodesInitCmd.MarkFlagRequired("image")

	nodesInitCmd.Flags().StringVar(&role, "role", "", "role for the node, i.e., peer, validator, nats")
	nodesInitCmd.MarkFlagRequired("role")

	requireInfrastructureFlags(nodesInitCmd, false)

	nodesInitCmd.Flags().StringVar(&healthCheckPath, "health-check-path", "", "path for the http health check on the node")
	nodesInitCmd.Flags().StringVar(&tcpIngressPorts, "tcp-ingress", "", "tcp ingress ports to open on the node")
	nodesInitCmd.Flags().StringVar(&udpIngressPorts, "udp-ingress", "", "udp ingress ports to open on the node")
	nodesInitCmd.Flags().StringVar(&taskRole, "task-role", "", "the optional vendor-specific task role (i.e., the ECS task execution role in the case of AWS)")
}
