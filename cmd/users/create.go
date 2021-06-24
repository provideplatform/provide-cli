package users

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

// createCmd represents the authenticate command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long:  `Create a new user in the configured ident instance; defaults to ident.provide.services.`,
	Run:   create,
}

var firstName string
var lastName string
var email string
var passwd string

func create(cmd *cobra.Command, args []string) {
	firstName = common.FreeInput("First Name", "", common.MandatoryValidation)
	lastName = common.FreeInput("Last Name", "", common.MandatoryValidation)
	email = common.FreeInput("Email", "", common.MandatoryValidation)
	passwd = common.FreeInput("Password", "", common.MandatoryValidation)

	resp, err := provide.CreateUser("", map[string]interface{}{
		"email":      email,
		"password":   passwd,
		"first_name": firstName,
		"last_name":  lastName,
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	_, err = provide.Authenticate(email, passwd)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Printf("created user: %s", resp.ID.String())
}
