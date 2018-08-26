package cmd

import (
	"fmt"
	"log"
	"os"

	provide "github.com/provideservices/provide-go"
	"github.com/spf13/cobra"
)

var walletsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of managed signing identities",
	Long:  `Retrieve a list of managed signing identities scoped to the authorized API token`,
	Run:   listWallets,
}

func listWallets(cmd *cobra.Command, args []string) {
	token := requireAPIToken()
	params := map[string]interface{}{}
	_, resp, err := provide.ListWallets(token, params)
	if err != nil {
		log.Printf("Failed to retrieve wallets list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp.([]interface{}) {
		wallet := resp.([]interface{})[i].(map[string]interface{})
		result := fmt.Sprintf("%s\t%s\n", wallet["id"], wallet["address"])
		if name, nameOk := wallet["name"].(string); nameOk {
			result = fmt.Sprintf("%s\t%s - %s\n", wallet["id"], name, wallet["address"])
		}
		fmt.Print(result)
	}
}

func init() {
	// networksListCmd.Flags().BoolVarP(&public, "public", "p", false, "filter private networks (false by default)")
}
