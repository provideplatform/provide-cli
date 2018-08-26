package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var apiTokensListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of API tokens",
	Long:  `Retrieve a list of API tokens scoped to the authorized API token`,
	Run:   listAPITokens,
}

func listAPITokens(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	_, resp, err := provide.ListTokens(token, params)
	if err != nil {
		log.Printf("Failed to retrieve API tokens list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp.([]interface{}) {
		dapp := resp.([]interface{})[i].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", dapp["id"], dapp["token"])
		fmt.Print(result)
	}
}
