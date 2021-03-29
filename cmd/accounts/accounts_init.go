package accounts

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	provide "github.com/provideservices/provide-go/api/nchain"
	providecrypto "github.com/provideservices/provide-go/crypto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountName string
var nonCustodial bool

var accountsInitCmd = &cobra.Command{
	Use:   "init [--non-custodial|-nc] [--network 024ff1ef-7369-4dee-969c-1918c6edb5d4]",
	Short: "Generate a new keypair for signing transactions and storing value",
	Long:  `Initialize a new account, which may be managed by Provide or you`,
	Run:   CreateAccount,
}

func CreateAccount(cmd *cobra.Command, args []string) {
	if nonCustodial {
		createDecentralizedAccount()
		return
	}

	createManagedAccount()
}

func createDecentralizedAccount() {
	publicKey, privateKey, err := providecrypto.EVMGenerateKeyPair()
	if err != nil {
		log.Printf("Failed to genereate nonCustodial keypair; %s", err.Error())
		os.Exit(1)
	}
	secret := hex.EncodeToString(providecrypto.FromECDSA(privateKey))
	keypairJSON, err := providecrypto.EVMMarshalEncryptedKey(providecrypto.HexToAddress(*publicKey), privateKey, secret)
	if err != nil {
		log.Printf("Failed to genereate nonCustodial keypair; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\n", *publicKey, string(keypairJSON))
	fmt.Print(result)
}

func createManagedAccount() {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"network_id": common.NetworkID,
	}
	if accountName != "" {
		params["name"] = accountName
	}
	account, err := provide.CreateAccount(token, params)
	if err != nil {
		log.Printf("Failed to genereate keypair; %s", err.Error())
		os.Exit(1)
	}

	common.AccountID = account.ID.String()
	result := fmt.Sprintf("Account %s\t%s\n", account.ID.String(), account.Address)
	// FIXME-- when account.Name exists... result = fmt.Sprintf("Account %s\t%s - %s\n", *account.Name, account.ID.String(), *account.Address)
	appAccountKey := common.BuildConfigKeyWithApp(common.AccountConfigKeyPartial, common.ApplicationID)
	if !viper.IsSet(appAccountKey) {
		viper.Set(appAccountKey, account.ID.String())
		viper.WriteConfig()
	}
	fmt.Print(result)
}

func init() {
	accountsInitCmd.Flags().BoolVarP(&nonCustodial, "non-custodial", "", false, "if the generated keypair is non-custodial")
	accountsInitCmd.Flags().StringVarP(&accountName, "name", "n", "", "human-readable name to associate with the generated keypair")
	accountsInitCmd.MarkFlagRequired("network")
}
