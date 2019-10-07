package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var applicationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of applications",
	Long:  `Retrieve a list of applications scoped to the authorized API token`,
	Run:   listApplications,
}

func listApplications(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	status, resp, err := provide.ListApplications(token, params)
	if err != nil {
		log.Printf("Failed to retrieve applications list; %s", err.Error())
		os.Exit(1)
	}
	if status != 200 {
		log.Printf("Failed to retrieve applications list; received status: %d", status)
		os.Exit(1)
	}
	for i := range resp.([]interface{}) {
		application := resp.([]interface{})[i].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", application["id"], application["name"])
		fmt.Print(result)
	}
}
