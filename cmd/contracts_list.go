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
	Short: "Retrieve a list of API tokens",
	Long:  `Retrieve a list of API tokens scoped to the authorized API token`,
	Run:   listAPITokens,
}

func listContracts(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	_, resp, err := provide.ListTokens(token, params)
	if err != nil {
		log.Printf("Failed to retrieve contracts list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp.([]interface{}) {
		contract := resp.([]interface{})[i].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", contract["id"], contract["name"])
		fmt.Print(result)
	}
}
