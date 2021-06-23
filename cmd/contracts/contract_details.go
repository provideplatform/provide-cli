package contracts

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var contractsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific smart contract",
	Long:  `Retrieve details for a specific smart contract by identifier, scoped to the authorized API token`,
	Run:   fetchContractDetails,
}

func fetchContractDetails(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	contract, err := provide.GetContractDetails(token, common.ContractID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for contract with id: %s; %s", common.ContractID, err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve details for contract with id: %s; %s", common.ContractID, resp)
	// 	os.Exit(1)
	// }
	result := fmt.Sprintf("%s\t%s\n", contract.ID.String(), *contract.Name)
	fmt.Print(result)
}

func init() {
	contractsDetailsCmd.Flags().StringVar(&common.ContractID, "contract", "", "id of the contract")
	contractsDetailsCmd.MarkFlagRequired("contract")
}
