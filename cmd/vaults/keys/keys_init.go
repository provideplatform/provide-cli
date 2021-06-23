package keys

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/cmd/common"
	vault "github.com/provideservices/provide-go/api/vault"

	"github.com/spf13/cobra"
)

var name string
var description string
var keyspec string
var keytype string
var keyusage string

var keysInitCmd = &cobra.Command{
	Use:   "init --name 'My Key' --description 'not your keys, not your crypto'",
	Short: "Create a new key",
	Long:  `Initialize a new key`,
	Run:   createKey,
}

func createKey(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func createKeyRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"name":        name,
		"description": description,
		"spec":        keyspec,
		"type":        keytype,
		"usage":       keyusage,
	}
	vlt, err := vault.CreateKey(token, common.VaultID, params)
	if err != nil {
		log.Printf("failed to create key in vault: %s; %s", common.VaultID, err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\t%s\n", vlt.ID.String(), *vlt.Name, *vlt.Description)
	fmt.Print(result)
}

func init() {
	keysInitCmd.Flags().StringVar(&name, "name", "", "name of the key")
	keysInitCmd.Flags().StringVar(&description, "description", "", "description of the key")
	keysInitCmd.Flags().StringVar(&keyspec, "spec", "", "key spec to use for the key")
	keysInitCmd.Flags().StringVar(&keytype, "type", "", "key type; must be symmetric or asymmetric")
	keysInitCmd.Flags().StringVar(&keyusage, "usage", "", "intended usage for the key; must be encrypt/decrypt or sign/verify")

	keysInitCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier for which the key will be created")
	keysInitCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization identifier for which the key will be created")
}
