package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// authenticateCmd represents the authenticate command
var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Authenticate using your developer credentials and receive a valid API token",
	Long: `Authenticate using user credentials retrieved from provide.services and receive a
valid API token which can be used to access the networks and application APIs.`,
	Run: authenticate,
}

func init() {
	rootCmd.AddCommand(authenticateCmd)
}

func authenticate(cmd *cobra.Command, args []string) {
	email := doEmailPrompt()
	passwd := doPasswordPrompt()

	fmt.Printf("Email/pw: %s :: %s", email, passwd)
	// TODO: IdentService.authenticate
	// TODO: if successful, write API token to secure configuration
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
	reader := bufio.NewReader(os.Stdin)
	passwd, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if passwd == "" {
		log.Println("Failed to read password from stdin")
		os.Exit(1)
	}
	return strings.Trim(passwd, "\n")
}
