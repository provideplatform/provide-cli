package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/provideservices/provide-go"
	"github.com/spf13/cobra"
)

var contractsDeployCmd = &cobra.Command{
	Use:   "deploy --application 8fec625c-a8ad-4197-bb77-8b46d7aecd8f --compile",
	Short: "Deploy a compiled smart contract on behalf of a dapp",
	Long:  `Deploy a smart contract on behalf of a dapp, optionally including ABI metadata for contracts referenced by CREATE opcodes`,
	Run:   deployContract,
}

func deployCompiledContract(bytecode []byte, constructorParams []interface{}) error {
	var err error
	keys := make([]int, 0)
	for k := range constructorParams {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var calldata string
	for _, k := range keys {
		calldata = fmt.Sprintf("%s%s", calldata, constructorParams[k])
	}
	return err
}

// deployContract deploys previously-compiled smart contract's ABI and bytecode along with any metadata
func deployContract(cmd *cobra.Command, args []string) {
	token := requireAPIToken()

	compileContract(nil, args)
	deployableArtifact, err := parseCachedArtifact()
	if err != nil {
		log.Printf("Failed to resolve deployable artifact; %s", err.Error())
		teardownAndExit(1)
	}

	params := map[string]interface{}{
		"application_id": applicationID,
		"data":           deployableArtifact["bytecode"],
		"lang":           "bytecode",
		"network_id":     networkID,
		"params":         deployableArtifact,
		"wallet_id":      walletID,
		"value":          0,
	}

	status, resp, err := provide.CreateTransaction(token, params)
	if err != nil {
		log.Printf("Failed to deploy contract; %s", err.Error())
		teardownAndExit(1)
	}
	if status == 201 {
		contract := resp.(map[string]interface{})
		fmt.Printf("Contract \t%s\n", contract["id"])
	} else {
		fmt.Printf("Failed to deploy contract; %s", resp)
		teardownAndExit(1)
	}
}

func init() {
	contractsDeployCmd.Flags().StringVar(&applicationID, "application", "", "application identifier this will belong to")
	contractsDeployCmd.MarkFlagRequired("application")

	contractsDeployCmd.Flags().StringVar(&networkID, "network", "", "network id (i.e., the network associated with the dapp)")
	contractsDeployCmd.MarkFlagRequired("network")

	contractsDeployCmd.Flags().StringVar(&walletID, "wallet", "", "wallet id which will be the signer of the contract creation tx")
	contractsDeployCmd.MarkFlagRequired("wallet")

	contractsDeployCmd.Flags().StringVar(&compilerVersion, "compiler-version", "latest", "target compiler version")
	contractsDeployCmd.Flags().StringVar(&compileWorkdir, "workdir", "", "path to temporary working directory for compiled artifacts")
	contractsDeployCmd.Flags().BoolVar(&skipOpcodesAnalysis, "skip-opcodes-analysis", false, "when true, static analysis of assembly for contract-internal ABI metadata is skipped")
	contractsDeployCmd.Flags().IntVar(&compilerOptimizerRuns, "optimizer-runs", 200, "set the number of runs to optimize for in terms of initial deployment cost; higher values optimize more for high-frequency usage; may not be supported by all compilers")
}
