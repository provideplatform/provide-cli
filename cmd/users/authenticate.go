package users

import (
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/ident"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// authenticateCmd represents the authenticate command
var AuthenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Authenticate using your credentials",
	Long: `Authenticate using user credentials and receive a
valid access/refresh token pair which can be used to make API calls.`,
	Run: authenticate,
}

func authenticate(cmd *cobra.Command, args []string) {
	RefreshToken()

	email = common.FreeInput("Email", "", common.MandatoryValidation)
	passwd = common.FreeInput("Password", "", common.MandatoryValidation)

	resp, err := provide.Authenticate(email, passwd)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if resp.Token.AccessToken != nil && resp.Token.RefreshToken != nil {
		cacheAccessRefreshToken(resp.Token)
	} else if resp.Token.Token != nil {
		cacheAPIToken(*resp.Token.Token)
	}

	log.Printf("Authentication successful")
}

func RefreshToken() {
	refreshToken := ""
	if viper.IsSet(common.RefreshTokenConfigKey) {
		refreshToken = viper.GetString(common.RefreshTokenConfigKey)
	}

	resp, err := provide.CreateToken(refreshToken, map[string]interface{}{
		"grant_type": "refresh_token",
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if resp != nil {
		cacheAccessRefreshToken(resp)
	}

	log.Printf("Refreshed access token")
}

func cacheAPIToken(token string) {
	viper.Set(common.AccessTokenConfigKey, token)
	viper.WriteConfig()
}

func cacheAccessRefreshToken(token *provide.Token) {
	if token.AccessToken != nil {
		viper.Set(common.AccessTokenConfigKey, *token.AccessToken)
	}

	if token.RefreshToken != nil {
		viper.Set(common.RefreshTokenConfigKey, *token.RefreshToken)
	}

	viper.WriteConfig()
}
