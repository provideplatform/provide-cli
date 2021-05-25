package accounts

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/nchain"
	"github.com/spf13/cobra"
)

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of signing identities",
	Long:  `Retrieve a list of signing identities (accounts) scoped to the authorized API token`,
	Run:   listAccounts,
}

func listAccounts(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	resp, err := provide.ListAccounts(token, params)
	if err != nil {
		log.Printf("Failed to retrieve accounts list; %s", err.Error())
		os.Exit(1)
	}
	// if status != 200 {
	// 	log.Printf("Failed to retrieve accounts list; received status: %d", status)
	// 	os.Exit(1)
	// }
	for i := range resp {
		account := resp[i]
		result := fmt.Sprintf("%s\t%s\n", account.ID.String(), account.Address)
		// TODO-- when account.Name exists... result = fmt.Sprintf("%s\t%s - %s\n", name, account, *account.Address)
		fmt.Print(result)
	}
}

func init() {
	accountsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter accounts")
	accountsListCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
