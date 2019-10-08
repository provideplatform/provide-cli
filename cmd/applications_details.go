package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var applicationsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific application",
	Long:  `Retrieve details for a specific application by identifier, scoped to the authorized API token`,
	Run:   fetchApplicationDetails,
}

func fetchApplicationDetails(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	status, resp, err := provide.GetApplicationDetails(token, applicationID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for application with id: %s; %s", applicationID, err.Error())
		os.Exit(1)
	}
	if status != 200 {
		log.Printf("Failed to retrieve details for application with id: %s; %s", applicationID, resp)
		os.Exit(1)
	}
	application = resp.(map[string]interface{})
	result := fmt.Sprintf("%s\t%s\n", application["id"], application["name"])
	fmt.Print(result)
}

func init() {
	applicationsDetailsCmd.Flags().StringVar(&applicationID, "application", "", "id of the application")
	applicationsDetailsCmd.MarkFlagRequired("application")
}
