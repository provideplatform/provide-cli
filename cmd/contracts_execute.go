package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var contractExecMethod string
var contractExecValue uint
var contractExecParams []interface{}

var contractsExecuteCmd = &cobra.Command{
	Use:   "execute --contract 0x5E250bB077ec836915155229E83d187715266167 --method vote --argv [] --value 0 --wallet 0x8A70B0C7E9896ac7025279a2Da240aEBD17A0cA3",
	Short: "Execute a smart contract",
	Long:  `Execute a smart contract method on a specific specific contract`,
	Run:   executeContract,
}

func executeContract(cmd *cobra.Command, args []string) {
	if walletID == "" {
		fmt.Println("Cannot execute a contract without a specified signer.")
		os.Exit(1)
	}
	token := requireAPIToken()
	params := map[string]interface{}{
		"method":    contractExecMethod,
		"params":    contractExecParams,
		"value":     contractExecValue,
		"wallet_id": walletID,
	}
	status, resp, err := provide.ExecuteContract(token, contractID, params)
	if err != nil {
		log.Printf("Failed to execute contract with id: %s; %s", contractID, err.Error())
		os.Exit(1)
	}
	if status == 200 {
		execution := resp.(map[string]interface{})
		executionJSON, _ := json.Marshal(execution)
		fmt.Printf("Successfully executed contract; response: %s", string(executionJSON))
	} else if status == 202 {
		execution := resp.(map[string]interface{})
		txRef := execution["ref"].(string)
		fmt.Printf("Successfully queued tx for asynchronous contract execution; tx ref: %s", txRef)
	} else if status >= 400 {
		fmt.Printf("Failed to execute contract; %d response: %s", status, resp)
	}
}

func init() {
	contractsExecuteCmd.Flags().StringVar(&contractID, "contract", "", "target contract id")
	contractsExecuteCmd.MarkFlagRequired("contract")

	contractsExecuteCmd.Flags().StringVar(&contractExecMethod, "method", "", "ABI method to invoke on the contract")
	contractsExecuteCmd.MarkFlagRequired("method")

	contractsExecuteCmd.Flags().UintVar(&contractExecValue, "value", 0, "value to send with transaction, specific in the smallest denonination of currency for the network (i.e., wei)")

	contractsExecuteCmd.Flags().StringVar(&walletID, "wallet", "", "wallet id with which to sign the tx")
	contractsExecuteCmd.MarkFlagRequired("wallet")
}
