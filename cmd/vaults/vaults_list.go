package vaults

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var pagination *common.Pagination

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
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
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
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(results))
	for i := range results {
		vlt := results[i]
		result := fmt.Sprintf("%s\t%s\t%s\n", vlt.ID.String(), *vlt.Name, *vlt.Description)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listVaultsRun,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	vaultsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter vaults")
	vaultsListCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier to filter vaults")
	vaultsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	vaultsListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	vaultsListCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	vaultsListCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of vaults to retrieve per page")
}
