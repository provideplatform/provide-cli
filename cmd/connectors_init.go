package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"
	"github.com/spf13/cobra"
)

// {"name": "IPFS us-east-1", "network_id": "aa51a87f-f142-4341-8e94-b4b0214a009f", "type": "ipfs", "config": {
// "region":"us-east-1",
// "target_id": "aws", "provider_id": "docker", "engine_id": "ipfs", "role": "ipfs", "container": "providenetwork-node", "api_port": 5001, "gateway_port": 8080,
// "env": {"CLIENT": "ipfs"}, "credentials": {"aws_access_key_id": "AKIASZW45XWQQDZIVFS4", "aws_secret_access_key": "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"},
// }}'

const connectorTypeIPFS = "ipfs"

var connectorName string
var connectorType string

var ipfsAPIPort uint
var ipfsGatewayPort uint

var connectorsInitCmd = &cobra.Command{
	Use:   "init --name 'my storage connector' --type ipfs --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new connector",
	Long:  `Initialize a new connector and orchestrate any related resources`,
	Run:   createConnector,
}

func securityConfigFactory() map[string]interface{} {
	if connectorType == connectorTypeIPFS {
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
	return nil
}

func connectorConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"credentials": infrastructureCredentialsConfigFactory(),
		"region":      region,
		"target_id":   targetID,
		"provider_id": providerID,
		"engine_id":   connectorType,
		"role":        connectorType,
		"container":   container,
		"env": map[string]interface{}{
			"CLIENT": connectorType,
		},
	}

	securityCfg := securityConfigFactory()
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

func createConnector(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{
		"name":       connectorName,
		"network_id": networkID,
		"type":       connectorType,
		"config":     connectorConfigFactory(),
	}
	status, resp, err := provide.CreateConnector(token, params)
	if err != nil {
		log.Printf("Failed to initialize connector; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		connector := resp.(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", connector["id"], connector["name"])
		fmt.Print(result)
	}
}

func init() {
	connectorsInitCmd.Flags().StringVar(&connectorName, "name", "", "name of the connector")
	connectorsInitCmd.MarkFlagRequired("name")

	connectorsInitCmd.Flags().StringVar(&connectorType, "type", "", "type of the connector")
	connectorsInitCmd.MarkFlagRequired("type")

	connectorsInitCmd.Flags().StringVar(&applicationID, "application", "", "application id")
	connectorsInitCmd.MarkFlagRequired("application")

	connectorsInitCmd.Flags().StringVar(&networkID, "network", "", "target network id")
	connectorsInitCmd.MarkFlagRequired("network")

	requireInfrastructureFlags(connectorsInitCmd)

	connectorsInitCmd.Flags().UintVar(&ipfsAPIPort, "ipfs_api_port", 5001, "tcp listen port for the ipfs api")
	connectorsInitCmd.Flags().UintVar(&ipfsGatewayPort, "ipfs_gateway_port", 8080, "tcp listen port for the ipfs gateway")
}
