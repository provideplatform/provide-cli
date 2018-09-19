package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var dappsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of dapps",
	Long:  `Retrieve a list of dapps scoped to the authorized API token`,
	Run:   listApplications,
}

func listApplications(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	_, resp, err := provide.ListApplications(token, params)
	if err != nil {
		log.Printf("Failed to retrieve dapps list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp.([]interface{}) {
		dapp := resp.([]interface{})[i].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", dapp["name"], dapp["id"])
		fmt.Print(result)
	}
}
