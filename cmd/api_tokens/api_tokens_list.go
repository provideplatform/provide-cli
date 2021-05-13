package api_tokens

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/ident"

	"github.com/spf13/cobra"
)

var apiTokensListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of API tokens",
	Long:  `Retrieve a list of API tokens scoped to the authorized API token`,
	Run:   listAPITokens,
}

func listAPITokens(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, "List")
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	resp, err := provide.ListTokens(token, params)
	if err != nil {
		log.Printf("Failed to retrieve API tokens list; %s", err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve API tokens list; received status: %d", status)
	// 	os.Exit(1)
	// }
	for i := range resp {
		apiToken := resp[i]
		result := fmt.Sprintf("%s\t%s\n", apiToken.ID.String(), *apiToken.Token)
		fmt.Print(result)
	}
}

func init() {
	apiTokensListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter API tokens")
}
