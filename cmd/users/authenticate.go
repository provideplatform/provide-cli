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
	Short: "Authenticate using your developer credentials",
	Long: `Authenticate using user credentials retrieved from provide.services and receive a
valid API token which can be used to access various APIs.`,
	Run: authenticate,
}

var email string
var passwd string

func authenticate(cmd *cobra.Command, args []string) {
	emailPrompt()
	passwordPrompt()

	resp, err := provide.Authenticate(email, passwd)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// if status != 201 {
	// 	log.Println("Authentication failed")
	// 	os.Exit(1)
	// }

	cacheAPIToken(*resp.Token.Token)
	log.Printf("Authentication successful")
}

func cacheAPIToken(token string) {
	viper.Set(common.AuthTokenConfigKey, token)
	viper.WriteConfig()
}
