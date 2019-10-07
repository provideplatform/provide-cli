package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var contractName string
var compiledArtifact map[string]interface{}

var contractsInitCmd = &cobra.Command{
	Use:   "init --name 'Registry' --network 024ff1ef-7369-4dee-969c-1918c6edb5d4",
	Short: "Initialize a new smart contract",
	Long:  `Initialize a new smart contract on behalf of a specific application; this operation may result in the contract being deployed`,
	Run:   createContract,
}

func compiledArtifactFactory() map[string]interface{} {
	if compiledArtifact != nil {
		return compiledArtifact
	}

	// TODO: support cli and soljson compiler
	return map[string]interface{}{
		"name":        nil,
		"abi":         nil,
		"assembly":    nil,
		"bytecode":    nil,
		"deps":        nil,
		"opcodes":     nil,
		"raw":         nil,
		"source":      nil,
		"fingerprint": nil,
	}
}

func contractParamsFactory() map[string]interface{} {
	return map[string]interface{}{
		"wallet_id":         walletID,
		"compiled_artifact": compiledArtifactFactory(),
	}
}

func createContract(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{
		"name":           contractName,
		"network_id":     networkID,
		"application_id": applicationID,
		"params":         contractParamsFactory(),
	}
	status, resp, err := provide.CreateContract(token, params)
	if err != nil {
		log.Printf("Failed to initialize application; %s", err.Error())
		os.Exit(1)
	}
	if status == 201 {
		contract = resp.(map[string]interface{})
		contractID = contract["id"].(string)
		result := fmt.Sprintf("%s\t%s\n", contract["id"], contract["name"])
		fmt.Print(result)
	}
	if !withoutAPIToken {
		createAPIToken(cmd, args)
	}
	if !withoutWallet {
		createWallet(cmd, args)
	}
}

func init() {
	contractsInitCmd.Flags().StringVar(&contractName, "name", "", "name of the contract")
	contractsInitCmd.MarkFlagRequired("name")

	contractsInitCmd.Flags().StringVar(&networkID, "network", "", "target network id")
	contractsInitCmd.MarkFlagRequired("network")

	contractsInitCmd.Flags().StringVar(&applicationID, "application", "", "target application id")
	contractsInitCmd.MarkFlagRequired("application")
}
