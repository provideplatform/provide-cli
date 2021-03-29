package api_tokens

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/ident"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiTokensInitCmd = &cobra.Command{
	Use:   "init --application 8fec625c-a8ad-4197-bb77-8b46d7aecd8f",
	Short: "Creates a new API token",
	Long:  `Initialize a new application API Token`,
	Run:   createAPIToken,
}

// createAPIToken triggers the generation of an API token for the given network.
func createAPIToken(cmd *cobra.Command, args []string) {
	token := common.RequireUserAuthToken()
	params := map[string]interface{}{}
	apiToken, err := provide.CreateApplicationToken(token, common.ApplicationID, params)
	if err != nil {
		log.Printf("Failed to create API token; %s", err.Error())
		os.Exit(1)
	}
	appAPITokenKey := common.BuildConfigKeyWithApp(common.APITokenConfigKeyPartial, common.ApplicationID)
	if !viper.IsSet(appAPITokenKey) {
		viper.Set(appAPITokenKey, apiToken.Token)
		viper.WriteConfig()
	}
	fmt.Printf("API Token\t%s\n", *apiToken.Token)
}

func init() {
	apiTokensInitCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application id")
	apiTokensInitCmd.MarkFlagRequired("application")
}
