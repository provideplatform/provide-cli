package vaults

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

var vaultsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of vaults",
	Long:  `Retrieve a list of vaults scoped to the authorized API token`,
	Run:   listVaults,
}

func listVaults(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listVaultsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	if common.OrganizationID != "" {
		params["organization_id"] = common.OrganizationID
	}
	results, resp, err := provide.ListVaults(token, params)
	if err != nil {
		log.Printf("failed to retrieve vaults list; %s", err.Error())
		os.Exit(1)
	}
	// TODO: better conditions logic here...
	firstRecordCount := 1
	if page > 1 {
		firstRecordCount = int(rpp) * int(page-1)
	}
	secondRecordCount := firstRecordCount + len(results)
	if firstRecordCount == 1 {
		secondRecordCount = secondRecordCount - 1
	}
	if len(results) != 0 {
		fmt.Printf("Showing record(s) %d-%d out of %s record(s)\n", firstRecordCount, secondRecordCount, resp.TotalCount)
	} else {
		fmt.Println("No more records found")
	}
	for i := range results {
		vlt := results[i]
		result := fmt.Sprintf("%s\t%s\t%s\n", vlt.ID.String(), *vlt.Name, *vlt.Description)
		fmt.Print(result)
	}
}

func init() {
	vaultsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter vaults")
	vaultsListCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier to filter vaults")
	vaultsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	vaultsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	vaultsListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	vaultsListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of vaults to retrieve per page")
}
