package users

import (
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AuthenticateCmd represents the authenticate command
var AuthenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Authenticate using your credentials",
	Long: `Authenticate using user credentials and receive a
valid access/refresh token pair which can be used to make API calls.`,
	Run: authenticate,
}

func authenticate(cmd *cobra.Command, args []string) {
	email = common.FreeInput("Email", "", common.MandatoryValidation)
	passwd = common.FreeInput("Password", "", common.MandatoryValidation)

	resp, err := provide.Authenticate(email, passwd)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if resp.Token.AccessToken != nil && resp.Token.RefreshToken != nil {
		common.CacheAccessRefreshToken(resp.Token)
	} else if resp.Token.Token != nil {
		cacheAPIToken(*resp.Token.Token)
	} else {
		log.Println("Failed to get token from authentication response.")
		os.Exit(1)
	}

	log.Print("Authentication successful")
	log.Printf("User ID: %s", common.DecodeUserID())
}

func cacheAPIToken(token string) {
	viper.Set(common.AccessTokenConfigKey, token)
	viper.WriteConfig()
}
