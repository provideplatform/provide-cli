package wallets

import (
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/nchain"

	"github.com/spf13/cobra"
)

var walletsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of custodial HD wallets",
	Long:  `Retrieve a list of HD wallets scoped to the authorized API token`,
	Run:   listWallets,
}

func listWallets(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepList)
}

func listWalletsRun(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	resp, err := provide.ListWallets(token, params)
	if err != nil {
		log.Printf("Failed to retrieve wallets list; %s", err.Error())
		os.Exit(1)
	}
	for i := range resp {
		wallet := resp[i]
		result := fmt.Sprintf("%s\t%s\n", wallet.ID.String(), *wallet.PublicKey)
		// FIXME-- when wallet.Name exists... result = fmt.Sprintf("Wallet %s\t%s - %s\n", wallet.Name, wallet.ID.String(), *wallet.Address)
		fmt.Print(result)
	}
}

func init() {
	walletsListCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application identifier to filter HD wallets")
}
