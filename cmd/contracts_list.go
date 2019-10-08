package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var contractsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of contracts",
	Long:  `Retrieve a list of contracts scoped to the authorized API token`,
	Run:   listContracts,
}

func listContracts(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	if applicationID != "" {
		params["application_id"] = applicationID
	}
	status, resp, err := provide.ListContracts(token, params)
	if err != nil {
		log.Printf("Failed to retrieve contracts list; %s", err.Error())
		os.Exit(1)
	}
	if status != 200 {
		log.Printf("Failed to retrieve contracts list; received status: %d", status)
		os.Exit(1)
	}
	contracts = resp.([]interface{})
	for i := range contracts {
		contract := resp.([]interface{})[i].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\t%s\n", contract["id"], contract["address"], contract["name"])
		fmt.Print(result)
	}
}

func init() {
	contractsListCmd.Flags().StringVar(&applicationID, "application", "", "application identifier to filter contracts")
}
