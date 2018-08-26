package cmd

import (
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
	log.Printf("Retrieved networks list:\n%s", resp)
}

func init() {
	networksListCmd.Flags().BoolVarP(&public, "public", "p", false, "filter private networks (false by default)")
}
