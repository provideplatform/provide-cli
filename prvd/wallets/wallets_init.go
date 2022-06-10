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

package wallets

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

var nonCustodial bool
var walletName string
var purpose int
var optional bool
var paginate bool

var walletsInitCmd = &cobra.Command{
	Use:   "init [--non-custodial|-nc]",
	Short: "Generate a new HD wallet for deterministically managing accounts, signing transactions and storing value",
	Long:  `Initialize a new HD wallet, which may be managed by Provide or you`,
	Run:   CreateWallet,
}

func CreateWallet(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func CreateWalletRun(cmd *cobra.Command, args []string) {
	if nonCustodial {
		createNonCustodialWallet()
		return
	}
	createCustodialWallet(cmd, args)
}

func createNonCustodialWallet() {
	publicKey, privateKey, err := providecrypto.EVMGenerateKeyPair()
	if err != nil {
		log.Printf("Failed to genereate non-custodial HD wallet; %s", err.Error())
		os.Exit(1)
	}
	secret := hex.EncodeToString(providecrypto.FromECDSA(privateKey))
	walletJSON, err := providecrypto.EVMMarshalEncryptedKey(providecrypto.HexToAddress(*publicKey), privateKey, secret)
	if err != nil {
		log.Printf("Failed to genereate non-custodial HD wallet; %s", err.Error())
		os.Exit(1)
	}
	result := fmt.Sprintf("%s\t%s\n", *publicKey, string(walletJSON))
	fmt.Print(result)
}

func createCustodialWallet(cmd *cobra.Command, args []string) {
	token := common.RequireAPIToken()
	params := map[string]interface{}{
		"purpose": purpose,
	}

	wallet, err := provide.CreateWallet(token, params)
	if err != nil {
		log.Printf("Failed to genereate custodial HD wallet; %s", err.Error())
		os.Exit(1)
	}
	common.WalletID = wallet.ID.String()
	result := fmt.Sprintf("Wallet %s\t%s\n", wallet.ID.String(), *wallet.PublicKey)

	if common.ApplicationID != "" {
		appWalletKey := common.BuildConfigKeyWithID(common.WalletConfigKeyPartial, common.ApplicationID)
		if !viper.IsSet(appWalletKey) {
			viper.Set(appWalletKey, wallet.ID.String())
			viper.WriteConfig()
		}
	} else if common.OrganizationID != "" {
		orgWalletKey := common.BuildConfigKeyWithID(common.WalletConfigKeyPartial, common.OrganizationID)
		if !viper.IsSet(orgWalletKey) {
			viper.Set(orgWalletKey, wallet.ID.String())
			viper.WriteConfig()
		}
	}

	fmt.Print(result)
}

func init() {
	walletsInitCmd.Flags().BoolVarP(&nonCustodial, "non-custodial", "", false, "if the generated HD wallet is custodial")
	walletsInitCmd.Flags().StringVarP(&walletName, "name", "n", "", "human-readable name to associate with the generated HD wallet")
	walletsInitCmd.Flags().IntVarP(&purpose, "purpose", "p", 44, "purpose of the HD wallet per BIP44")
	walletsInitCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
