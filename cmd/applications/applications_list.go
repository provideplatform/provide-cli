package applications

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var applicationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of applications",
	Long:  `Retrieve a list of applications scoped to the authorized API token`,
	Run:   listApplications,
}

func listApplications(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	applications, err := provide.ListApplications(token, params)
	if err != nil {
		log.Printf("Failed to retrieve applications list; %s", err.Error())
		os.Exit(1)
	}
	for i := range applications {
		application := applications[i]
		result := fmt.Sprintf("%s\t%s\n", application.ID.String(), *application.Name)
		fmt.Print(result)
	}
}
