package vaults

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var vaultsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of vaults",
	Long:  `Retrieve a list of vaults scoped to the authorized API token`,
	Run:   listVaults,
}

func listVaults(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listVaultsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	if common.OrganizationID != "" {
		params["organization_id"] = common.OrganizationID
	}
	resp, err := provide.ListVaults(token, params)
	if err != nil {
		log.Printf("failed to retrieve vaults list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp {
		vlt := resp[i]
		result := fmt.Sprintf("%s\t%s\t%s\n", vlt.ID.String(), *vlt.Name, *vlt.Description)
		fmt.Print(result)
	}
}

func init() {
	vaultsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter vaults")
	vaultsListCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier to filter vaults")
	vaultsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
