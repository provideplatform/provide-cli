package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var contractsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific smart contract",
	Long:  `Retrieve details for a specific smart contract by identifier, scoped to the authorized API token`,
	Run:   fetchContractDetails,
}

func fetchContractDetails(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	status, resp, err := provide.GetContractDetails(token, contractID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for contract with id: %s; %s", contractID, err.Error())
		os.Exit(1)
	}
	if status != 200 {
		log.Printf("Failed to retrieve details for contract with id: %s; %s", contractID, resp)
		os.Exit(1)
	}
	contract = resp.(map[string]interface{})
	result := fmt.Sprintf("%s\t%s\n", contract["id"], contract["name"])
	fmt.Print(result)
}

func init() {
	contractsDetailsCmd.Flags().StringVar(&contractID, "contract", "", "id of the contract")
	contractsDetailsCmd.MarkFlagRequired("contract")
}
