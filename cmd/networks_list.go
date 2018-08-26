package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-go"

	"github.com/spf13/cobra"
)

var public bool

var networksListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of networks",
	Long:  `Retrieve a list of networks scoped to the authorized API token`,
	Run:   listNetworks,
}

func listNetworks(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	if public {
		params["public"] = "true"
	}
	_, resp, err := provide.ListNetworks(token, params)
	if err != nil {
		log.Printf("Failed to retrieve networks list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp.([]interface{}) {
		network := resp.([]interface{})[i].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", network["id"], network["name"])
		fmt.Print(result)
	}
}

func init() {
	networksListCmd.Flags().BoolVarP(&public, "public", "p", false, "filter private networks (false by default)")
}
