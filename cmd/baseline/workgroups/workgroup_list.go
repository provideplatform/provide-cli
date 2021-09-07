package workgroups

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	ident "github.com/provideplatform/provide-go/api/ident"
	"github.com/spf13/cobra"
)

var page uint64
var rpp uint64

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
	applications, err := ident.ListApplications(token, map[string]interface{}{
		"type": "baseline",
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	})
	if err != nil {
		log.Printf("failed to retrieve baseline workgroups; %s", err.Error())
		os.Exit(1)
	}
	for i := range applications {
		workgroup := applications[i]
		result := fmt.Sprintf("%s\t%s\n", workgroup.ID.String(), *workgroup.Name)
		fmt.Print(result)
	}
}

func init() {
	listBaselineWorkgroupsCmd.Flags().Uint64Var(&page, "page", 1, "page number to retrieve")
	listBaselineWorkgroupsCmd.Flags().Uint64Var(&rpp, "rpp", 25, "number of baseline workgroups to retrieve per page")
}
