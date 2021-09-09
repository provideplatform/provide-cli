package applications

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var applicationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of applications",
	Long:  `Retrieve a list of applications scoped to the authorized API token`,
	Run:   listApplications,
}

func listApplications(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	applications, _, err := provide.ListApplications(token, params)
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

func init() {
	applicationsListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	applicationsListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of applications to retrieve per page")
	applicationsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
