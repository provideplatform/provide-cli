package workgroups

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	ident "github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var pagination *common.Pagination

var listBaselineWorkgroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List baseline workgroups",
	Long:  `List all available baseline workgroups`,
	Run:   listWorkgroups,
}

func listWorkgroups(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listWorkgroupsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	applications, resp, err := ident.ListApplications(token, map[string]interface{}{
		"type": "baseline",
		"page": fmt.Sprintf("%d", pagination.Page),
		"rpp":  fmt.Sprintf("%d", pagination.Rpp),
	})
	if err != nil {
		log.Printf("failed to retrieve baseline workgroups; %s", err.Error())
		os.Exit(1)
	}
	pagination.UpdateCountsAndPrintCurrentInterval(resp.TotalCount, len(applications))
	for i := range applications {
		workgroup := applications[i]
		result := fmt.Sprintf("%s\t%s\n", workgroup.ID.String(), *workgroup.Name)
		fmt.Print(result)
	}

	paginationPrompt := &common.PaginationPrompt{
		Pagination:  pagination,
		CurrentStep: "",
		RunPageCmd:  listWorkgroupsRun,
	}
	common.AutoPromptPagination(cmd, args, paginationPrompt)
}

func init() {
	pagination = &common.Pagination{}
	listBaselineWorkgroupsCmd.Flags().IntVar(&pagination.Page, "page", common.DefaultPage, "page number to retrieve")
	listBaselineWorkgroupsCmd.Flags().IntVar(&pagination.Rpp, "rpp", common.DefaultRpp, "number of baseline workgroups to retrieve per page")
	listBaselineWorkgroupsCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
}
