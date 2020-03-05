package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var messageBusInitCmd = &cobra.Command{
	Use:   "init --name 'my magic bus app' --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new message bus application",
	Long:  `Initialize a new message bus application, its initial IPFS connector and on-chain registry contract`,
	Run:   createMessageBus,
}

func createMessageBus(cmd *cobra.Command, args []string) {
	createApplication(cmd, args)
	if application == nil {
		fmt.Println("Cannot continue provisioning message bus application without a valid application context.")
		os.Exit(1)
	}

	if connectorName == "" {
		connectorName = fmt.Sprintf("%s message bus %s connector - %s", application["name"], connectorType, region)
	}
	createConnector(cmd, args)
	if connector == nil {
		fmt.Println("Failed to provision connector for message bus application.")
		os.Exit(1)
	}

	createContract(cmd, args)
	if contract == nil {
		fmt.Println("Failed to deploy registry contract for message bus application.")
		os.Exit(1)
	}
}

func init() {
	applicationType = applicationTypeMessageBus
	initMessageBusRegistryContractCompiledArtifact()

	messageBusInitCmd.Flags().StringVar(&networkID, "network", "", "target network id")
	messageBusInitCmd.MarkFlagRequired("network")

	// application
	messageBusInitCmd.Flags().StringVar(&applicationName, "name", "", "name of the message bus application")
	messageBusInitCmd.MarkFlagRequired("name")

	// connector
	messageBusInitCmd.Flags().StringVar(&connectorName, "connector-name", "", "name of the connector")
	messageBusInitCmd.Flags().StringVar(&connectorType, "connector-type", "ipfs", "type of the connector")

	requireInfrastructureFlags(messageBusInitCmd, true)

	messageBusInitCmd.Flags().UintVar(&ipfsAPIPort, "connector-ipfs-api-port", 5001, "tcp listen port for the ipfs api")
	messageBusInitCmd.Flags().UintVar(&ipfsGatewayPort, "connector-ipfs-gateway-port", 8080, "tcp listen port for the ipfs gateway")

	// registry contract
	messageBusInitCmd.Flags().StringVar(&contractName, "contract-name", "Registry", "name of the registry contract")
	messageBusInitCmd.Flags().StringVar(&walletID, "wallet", "", "wallet id with which to sign the tx")
}
