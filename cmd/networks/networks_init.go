package networks

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/nchain"
	"github.com/spf13/cobra"
)

var chain string
var nativeCurrency string
var platform string
var protocolID string

var networkName string
var networksInitCmd = &cobra.Command{
	Use:   "init --name 'whiteblock testnet",
	Short: "Initialize a new network",
	Long:  `Initialize a new network with options`,
	Run:   CreateNetwork,
}

// CreateNetwork configures a new peer-to-peer network;
// see https://docs.provide.services/microservices/goldmine/#create-a-network
func CreateNetwork(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"name":   networkName,
		"config": configFactory(),
	}
	network, err := provide.CreateNetwork(token, params)
	if err != nil {
		log.Printf("Failed to initialize network; %s", err.Error())
		os.Exit(1)
	}
	common.NetworkID = network.ID.String()
	result := fmt.Sprintf("%s\t%s\n", network.ID.String(), *network.Name)
	fmt.Print(result)
}

func init() {
	networksInitCmd.Flags().StringVar(&networkName, "name", "", "name of the network")
	networksInitCmd.MarkFlagRequired("name")

	networksInitCmd.Flags().StringVar(&chain, "chain", "", "name of the chain")
	networksInitCmd.MarkFlagRequired("chain")

	networksInitCmd.Flags().StringVar(&common.EngineID, "engine", "", "consensus engine to be used for the chain (i.e., ethash, poa, ibft)")
	networksInitCmd.MarkFlagRequired("engine")

	networksInitCmd.Flags().StringVar(&nativeCurrency, "native-currency", "", "symbol representing the native currency on the network (i.e., ETH)")
	networksInitCmd.MarkFlagRequired("native-currency")

	networksInitCmd.Flags().StringVar(&platform, "platform", "", "platform type (i.e., evm, bcoin)")
	networksInitCmd.MarkFlagRequired("platform")

	networksInitCmd.Flags().StringVar(&protocolID, "protocol", "", "type of consensus mechanism (i.e., pow, poa)")
	networksInitCmd.MarkFlagRequired("protocol")

}

func configFactory() map[string]interface{} {
	var chainspec map[string]interface{}
	if common.EngineID == "clique" {
		chainspec = cliqueChainspecFactory()
	} else {
		log.Printf("Failed to initialize network; additional chainspec factories should be implemented")
		os.Exit(1)
	}

	return map[string]interface{}{
		"chain":           chain,
		"chainspec":       chainspec,
		"engine_id":       common.EngineID,
		"native_currency": nativeCurrency,
		"platform":        platform,
		"protocol_id":     protocolID,
	}
}

func cliqueChainspecFactory() map[string]interface{} {
	genesis := map[string]interface{}{}
	json.Unmarshal([]byte(`{
		"config": {
		  "chainId": <arbitrary positive integer>,
		  "homesteadBlock": 0,
		  "eip150Block": 0,
		  "eip155Block": 0,
		  "eip158Block": 0,
		  "byzantiumBlock": 0,
		  "constantinopleBlock": 0,
		  "petersburgBlock": 0
		},
		"alloc": {},
		"coinbase": "0x0000000000000000000000000000000000000000",
		"difficulty": "0x20000",
		"extraData": "",
		"gasLimit": "0x2fefd8",
		"nonce": "0x0000000000000042",
		"mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"timestamp": "0x00"
	}`), &genesis)
	return genesis
}
