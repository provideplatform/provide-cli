package info

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const contractTypeRegistry = "registry"

var contract map[string]interface{}
var contracts []interface{}
var contractType string

// InfoCmd is the handler for the `info` command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about the currently authorized user",
	Long:  "Get information about the currently authorized user",
	Run:   showInfo,
}

func showInfo(cmd *cobra.Command, args []string) {
	// User ID, email, username, ..

	user, err := common.GetUserDetails()

	if err != nil {
		log.Printf("Unable to get user details: %s", err)
		os.Exit(1)
	}

	if user != nil {
		fmt.Println("Current User:")
		fmt.Println(" ID:         ", user.ID)
		fmt.Println(" Name:       ", user.Name)
		fmt.Println(" First Name: ", user.FirstName)
		fmt.Println(" Last Name:  ", user.LastName)
		fmt.Println(" Email:      ", user.Email)
	}
}
