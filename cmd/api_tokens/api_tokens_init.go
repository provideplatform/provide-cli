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

var scope string
var grantType string
var offlineAccess bool
var refreshToken bool

var apiTokensInitCmd = &cobra.Command{
	Use:   "init [--application 8fec625c-a8ad-4197-bb77-8b46d7aecd8f] [--organization 2209cf15-2402-4e25-b6b6-1c901b9dde69] [--offline-access] [--refresh-token]",
	Short: "Authorize a new API access or refresh token",
	Long:  `Authorize a new API token on behalf of the given application or organization`,
	Run:   createAPIToken,
}

// createAPIToken triggers the generation of an API token for the given network.
func createAPIToken(cmd *cobra.Command, args []string) {
	userToken := common.RequireUserAuthToken()
	params := map[string]interface{}{}

	if scope != "" {
		params["scope"] = scope
	} else if offlineAccess {
		params["scope"] = "offline_access"
	}

	if grantType != "" {
		params["grant_type"] = grantType
	} else if refreshToken {
		params["grant_type"] = "refresh_token"
	}

	if common.ApplicationID != "" {
		token, err := provide.CreateApplicationToken(userToken, common.ApplicationID, params)
		if err != nil {
			log.Printf("Failed to authorize API token on behalf of application %s; %s", common.ApplicationID, err.Error())
			os.Exit(1)
		}

		appAPITokenKey := common.BuildConfigKeyWithApp(common.APITokenConfigKeyPartial, common.ApplicationID)
		if !viper.IsSet(appAPITokenKey) {
			viper.Set(appAPITokenKey, token.Token)
			viper.WriteConfig()
		}

		if token.Token != nil {
			fmt.Printf("API token authorized for application: %s\t%s\n", common.ApplicationID, *token.AccessToken)
		} else if token.AccessToken != nil {
			fmt.Printf("Access token authorized for application: %s\t%s\n", common.ApplicationID, *token.AccessToken)
			if token.RefreshToken != nil {
				fmt.Printf("Refresh token authorized for application: %s\t%s\n", common.ApplicationID, *token.RefreshToken)
			}
		}

	} else if common.OrganizationID != "" {
		token, err := provide.CreateApplicationToken(userToken, common.OrganizationID, params)
		if err != nil {
			log.Printf("Failed to authorize API token on behalf of organization %s; %s", common.ApplicationID, err.Error())
			os.Exit(1)
		}

		orgAPITokenKey := common.BuildConfigKeyWithOrg(common.APITokenConfigKeyPartial, common.OrganizationID)
		if !viper.IsSet(orgAPITokenKey) {
			viper.Set(orgAPITokenKey, token.Token)
			viper.WriteConfig()
		}

		if token.Token != nil {
			fmt.Printf("API token authorized for organization: %s\t%s\n", common.OrganizationID, *token.Token)
		} else if token.AccessToken != nil {
			fmt.Printf("Access token authorized for organization: %s\t%s\n", common.OrganizationID, *token.AccessToken)
			if token.RefreshToken != nil {
				fmt.Printf("Refresh token authorized for organization: %s\t%s\n", common.OrganizationID, *token.RefreshToken)
			}
		}
	}
}

func init() {
	apiTokensInitCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application id")
	apiTokensInitCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization id")

	apiTokensInitCmd.Flags().BoolVar(&offlineAccess, "offline-access", false, "offline access")
	apiTokensInitCmd.Flags().BoolVar(&refreshToken, "refresh-token", false, "refresh token")
}
