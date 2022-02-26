package networks

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var public bool

var page uint64
var rpp uint64

var networksListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of networks",
	Long:  `Retrieve a list of networks scoped to the authorized API token`,
	Run:   listNetworks,
}

func listNetworks(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"page": fmt.Sprintf("%d", page),
		"rpp":  fmt.Sprintf("%d", rpp),
	}
	if public {
		params["public"] = "true"
	}
	networks, err := provide.ListNetworks(token, params)
	if err != nil {
		log.Printf("Failed to retrieve networks list; %s", err.Error())
		os.Exit(1)
	}
	for i := range networks {
		network := networks[i]
		result := fmt.Sprintf("%s\t%s\n", network.ID.String(), *network.Name)
		fmt.Print(result)
	}
}

func init() {
	networksListCmd.Flags().BoolVarP(&public, "public", "p", false, "filter private networks (false by default)")
	networksListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
	networksListCmd.Flags().BoolVarP(&paginate, "paginate", "", false, "List pagination flags")
	networksListCmd.Flags().Uint64Var(&page, "page", common.DefaultPage, "page number to retrieve")
	networksListCmd.Flags().Uint64Var(&rpp, "rpp", common.DefaultRpp, "number of networks to retrieve per page")
}
