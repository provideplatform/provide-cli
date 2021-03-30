package users

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/ident"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// authenticateCmd represents the authenticate command
var AuthenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Authenticate using your developer credentials and receive a valid API token",
	Long: `Authenticate using user credentials retrieved from provide.services and receive a
valid API token which can be used to access the networks and application APIs.`,
	Run: authenticate,
}

func authenticate(cmd *cobra.Command, args []string) {
	email := doEmailPrompt()
	passwd := doPasswordPrompt()

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

func doEmailPrompt() string {
	fmt.Print("Email: ")
	reader := bufio.NewReader(os.Stdin)
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if email == "" {
		log.Println("Failed to read email from stdin")
		os.Exit(1)
	}
	return strings.Trim(email, "\n")
}

func doPasswordPrompt() string {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	passwd := string(password[:])
	if passwd == "" {
		log.Println("Failed to read password from stdin")
		os.Exit(1)
	}
	return strings.Trim(passwd, "\n")
}

func cacheAPIToken(token string) {
	viper.Set(common.AuthTokenConfigKey, token)
	viper.WriteConfig()
}
