package vaults

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var name string
var description string

var vaultsInitCmd = &cobra.Command{
	Use:   "init --name 'My Vault' --description 'not your keys, not your crypto'",
	Short: "Create a new vault",
	Long:  `Initialize a new vault`,
	Run:   createVault,
}

func createVault(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"name":        name,
		"description": description,
	}
	vlt, err := provide.CreateVault(token, params)
	if err != nil {
		log.Printf("Failed to genereate HD wallet; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\t%s\n", vlt.ID.String(), *vlt.Name, *vlt.Description)
	fmt.Print(result)
}

func init() {
	vaultsInitCmd.Flags().StringVar(&name, "name", "", "name of the vault")
	vaultsInitCmd.Flags().StringVar(&description, "description", "", "description of the vault")

	vaultsInitCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier for which the vault will be created")
	vaultsInitCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier for which the vault will be created")
}
