package workgroups

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	ident "github.com/provideservices/provide-go/api/ident"
	"github.com/spf13/cobra"
)

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

}
