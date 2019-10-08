package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/spf13/cobra"
)

var messageBusPublishCmd = &cobra.Command{
	Use:   "publish --application 57f8c0af-089a-4ab3-b3c2-ca8a9ed547e0 --wallet 0xEA490AA70a95D5dAcD9Cf6a141692847455Ba928 /path/to/file.json",
	Short: "Publish a message to a message bus application",
	Long: `Publish a message by writing it to the specified message bus application; the named input file is hashed
and its raw bytes are written to the configured IPFS connector, and its hash is published to the on-chain registry contract`,
	Run: publishMessage,
}

func publishMessage(cmd *cobra.Command, args []string) {
	fetchApplicationDetails(cmd, args)
	if application == nil {
		log.Printf("Failed to retrieve message bus application with id: %s", applicationID)
		os.Exit(1)
	}

	config, configOk := application["config"].(map[string]interface{})
	if !configOk {
		log.Printf("Failed to parse message bus application config for message bus application with id: %s", applicationID)
		os.Exit(1)
	}

	if applicationType, applicationTypeOk := config["type"].(string); applicationTypeOk && applicationType != applicationTypeMessageBus {
		log.Printf("Retrieved application with id %s, but it was not a message bus application", applicationID)
		os.Exit(1)
	}

	resolveMessageBusContract(cmd, args)
	resolveMessageBusConnector(cmd, args)

	connectorCfg := connector["config"].(map[string]interface{})
	connectorAPIURL, connectorAPIURLOk := connectorCfg["api_url"].(string)
	if !connectorAPIURLOk {
		log.Printf("No connector API URL resolved for message bus application with id: %s", applicationID)
		os.Exit(1)
	}
	// connectorGatewayURL, connectorGatewayURLOk := connectorCfg["gateway_url"].(string)=
	// if !connectorGatewayURLOk {
	// 	log.Printf("No connector gateway URL resolved for message bus application with id: %s", applicationID)
	// 	os.Exit(1)
	// }

	if len(args) == 0 {
		log.Printf("No message file path provided for publishing via message bus application: %s", applicationID)
		os.Exit(1)
	}

	file := args[0]
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Failed to read message file at path %s for publishing via message bus application: %s; %s", file, applicationID, err.Error())
		os.Exit(1)
	}

	sh := ipfs.NewShell(connectorAPIURL)
	hash, err := sh.Add(strings.NewReader(string(data)))
	if err != nil {
		log.Printf("Failed to publish %d-byte message via message bus application: %s; %s", len(data), applicationID, err.Error())
		os.Exit(1)
	}
	log.Printf("Published %d byte(s) to IPFS; hash: %s", len(data), hash)

	// TODO: resolve contract id
	contractExecParams = []interface{}{hash}
	executeContract(cmd, args)
	if err != nil {
		log.Printf("Failed to execute publish method on message bus application registry contract with id: %s; %s", contractID, err.Error())
		os.Exit(1)
	}
}

func resolveMessageBusContract(cmd *cobra.Command, args []string) {
	listContracts(cmd, args)
	if contracts == nil || len(contracts) == 0 {
		log.Printf("No contracts resolved for message bus application with id: %s", applicationID)
		os.Exit(1)
	}
	for _, c := range contracts {
		contractID = c.(map[string]interface{})["id"].(string)
		fetchContractDetails(cmd, args)
		if contract != nil {
			if params, paramsOk := contract["params"].(map[string]interface{}); paramsOk {
				if params["type"] == contractTypeRegistry {
					contractID = contract["id"].(string)
					break
				}
			}
			contract = nil
			contractID = ""
		}
	}

	if contract == nil {
		log.Printf("No registry contract resolved for message bus application with id %s", applicationID)
		os.Exit(1)
	}
}

func resolveMessageBusConnector(cmd *cobra.Command, args []string) {
	listConnectors(cmd, args)
	if connectors == nil || len(connectors) == 0 {
		log.Printf("No connectors resolved for message bus application with id: %s", applicationID)
		os.Exit(1)
	}
	for _, c := range connectors {
		cnnector := c.(map[string]interface{})
		if cnnector["type"] == connectorTypeIPFS {
			connector = cnnector
			break
		}
	}

	if connector == nil {
		log.Printf("No IPFS connector resolved for message bus application with id %s", applicationID)
		os.Exit(1)
	}
}

func init() {
	messageBusPublishCmd.Flags().StringVar(&applicationID, "application", "", "target message bus application id")
	messageBusPublishCmd.MarkFlagRequired("application")

	messageBusPublishCmd.Flags().StringVar(&walletID, "wallet", "", "id or address of the signer for the registry transaction")
	messageBusPublishCmd.MarkFlagRequired("wallet")
}
