package cmd

import (
	"fmt"
	"log"
	"os"

	provide "github.com/provideservices/provide-go"
	"github.com/spf13/cobra"
)

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of signing identities",
	Long:  `Retrieve a list of signing identities (accounts) scoped to the authorized API token`,
	Run:   listAccounts,
}

func listAccounts(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	if applicationID != "" {
		params["application_id"] = applicationID
	}
	status, resp, err := provide.ListAccounts(token, params)
	if err != nil {
		log.Printf("Failed to retrieve accounts list; %s", err.Error())
		os.Exit(1)
	}
	if status != 200 {
		log.Printf("Failed to retrieve accounts list; received status: %d", status)
		os.Exit(1)
	}
	for i := range resp.([]interface{}) {
		account := resp.([]interface{})[i].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", account["id"], account["address"])
		if name, nameOk := account["name"].(string); nameOk {
			result = fmt.Sprintf("%s\t%s - %s\n", name, account["id"], account["address"])
		}
		fmt.Print(result)
	}
}

func init() {
	accountsListCmd.Flags().StringVar(&applicationID, "application", "", "application identifier to filter accounts")
}
