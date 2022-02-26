package applications

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var applicationsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Retrieve a specific application",
	Long:  `Retrieve details for a specific application by identifier, scoped to the authorized API token`,
	Run:   fetchApplicationDetails,
}

func fetchApplicationDetails(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	application, err := provide.GetApplicationDetails(token, common.ApplicationID, params)
	if err != nil {
		log.Printf("Failed to retrieve details for application with id: %s; %s", common.ApplicationID, err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\n", application.ID.String(), *application.Name)
	fmt.Print(result)
}

func init() {
	applicationsDetailsCmd.Flags().StringVar(&common.ApplicationID, "application", "", "id of the application")
	// applicationsDetailsCmd.MarkFlagRequired("application")
}
