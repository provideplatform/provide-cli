package users

import (
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

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
	email = common.FreeInput("Email", "", common.MandatoryValidation)
	passwd = common.FreeInput("Password", "", common.MandatoryValidation)

	resp, err := provide.Authenticate(email, passwd)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if resp.Token.AccessToken != nil && resp.Token.RefreshToken != nil {
		common.CacheAccessRefreshToken(resp.Token, nil)
	} else if resp.Token.Token != nil {
		viper.Set(common.AccessTokenConfigKey, *resp.Token.Token)
		viper.WriteConfig()
	}

	log.Printf("Authentication successful")
}
