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

var connectorName string
var connectorType string

var ipfsAPIPort uint64
var ipfsGatewayPort uint64

var optional bool
var paginate bool

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
					"tcp": []uint64{ipfsAPIPort, ipfsGatewayPort},
					"udp": []uint{},
				},
			},
		}
	}
	return nil
}

func connectorConfigFactory() map[string]interface{} {
	cfg := map[string]interface{}{
		"credentials":   common.InfrastructureCredentialsConfigFactory(),
		"common.Image":  common.Image,
		"common.Region": common.Region,
		"target_id":     common.TargetID,
		"provider_id":   common.ProviderID,
		"engine_id":     connectorType,
		"role":          connectorType,
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
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"name":       connectorName,
		"network_id": common.NetworkID,
		"type":       connectorType,
		"config":     connectorConfigFactory(),
	}
	connector, err := provide.CreateConnector(token, params)
	if err != nil {
		log.Printf("Failed to initialize connector; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\n", connector.ID.String(), *connector.Name)
	fmt.Print(result)
}

func init() {
	connectorsInitCmd.Flags().StringVar(&connectorName, "name", "", "name of the connector")
	//connectorsInitCmd.MarkFlagRequired("name")

	connectorsInitCmd.Flags().StringVar(&connectorType, "type", "", "type of the connector")
	// connectorsInitCmd.MarkFlagRequired("type")

	connectorsInitCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application id")
	// connectorsInitCmd.MarkFlagRequired("application")

	connectorsInitCmd.Flags().StringVar(&common.NetworkID, "network", "", "target network id")
	// connectorsInitCmd.MarkFlagRequired("network")

	common.RequireInfrastructureFlags(connectorsInitCmd, true)

	connectorsInitCmd.Flags().Uint64Var(&ipfsAPIPort, "ipfs-api-port", 5001, "tcp listen port for the ipfs api")
	connectorsInitCmd.Flags().Uint64Var(&ipfsGatewayPort, "ipfs-gateway-port", 8080, "tcp listen port for the ipfs gateway")
	connectorsInitCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")

}
