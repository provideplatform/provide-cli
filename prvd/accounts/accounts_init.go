/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package accounts

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/nchain"
	providecrypto "github.com/provideplatform/provide-go/crypto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountName string
var nonCustodial bool
var optional bool
var paginate bool

var accountsInitCmd = &cobra.Command{
	Use:   "init [--non-custodial|-nc] [--network 024ff1ef-7369-4dee-969c-1918c6edb5d4] [--application 024ff1ef-7369-4dee-969c-1918c6edb5d4] [--organization 024ff1ef-7369-4dee-969c-1918c6edb5d4]",
	Short: "Generate a new keypair for signing transactions and storing value",
	Long:  `Initialize a new account, which may be managed by Provide or you`,
	Run:   CreateAccount,
}

func CreateAccount(cmd *cobra.Command, args []string) {
	if nonCustodial {
		createNonCustodialAccount()
		return
	}

	createManagedAccount(cmd, args)
}

func createNonCustodialAccount() {
	publicKey, privateKey, err := providecrypto.EVMGenerateKeyPair()
	if err != nil {
		log.Printf("Failed to genereate non-custodial keypair; %s", err.Error())
		os.Exit(1)
	}
	secret := hex.EncodeToString(providecrypto.FromECDSA(privateKey))
	keypairJSON, err := providecrypto.EVMMarshalEncryptedKey(providecrypto.HexToAddress(*publicKey), privateKey, secret)
	if err != nil {
		log.Printf("Failed to genereate non-custodial keypair; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\n", *publicKey, string(keypairJSON))
	fmt.Print(result)
}

func createManagedAccount(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"network_id": common.NetworkID,
	}
	if accountName != "" {
		params["name"] = accountName
	}
	if common.ApplicationID != "" {
		params["application_id"] = common.ApplicationID
	}
	if common.OrganizationID != "" {
		params["organization_id"] = common.OrganizationID
	}
	account, err := provide.CreateAccount(token, params)
	if err != nil {
		log.Printf("Failed to genereate keypair; %s", err.Error())
		os.Exit(1)
	}

	common.AccountID = account.ID.String()
	result := fmt.Sprintf("Account %s\t%s\n", account.ID.String(), account.Address)
	// FIXME-- when account.Name exists... result = fmt.Sprintf("Account %s\t%s - %s\n", *account.Name, account.ID.String(), *account.Address)
	appAccountKey := common.BuildConfigKeyWithID(common.AccountConfigKeyPartial, common.ApplicationID)
	if !viper.IsSet(appAccountKey) {
		viper.Set(appAccountKey, account.ID.String())
		viper.WriteConfig()
	}
	fmt.Print(result)
}

func init() {
	accountsInitCmd.Flags().BoolVarP(&nonCustodial, "non-custodial", "", false, "if the generated keypair is non-custodial")
	accountsInitCmd.Flags().StringVarP(&accountName, "name", "n", "", "human-readable name to associate with the generated keypair")

	accountsInitCmd.Flags().StringVar(&common.NetworkID, "network", "", "network id")
	accountsInitCmd.MarkFlagRequired("network")

	accountsInitCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application id")
	accountsInitCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization id")
	accountsInitCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
