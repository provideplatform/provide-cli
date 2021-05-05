package workgroups

import (
	"log"
	"os"

	"github.com/provideservices/provide-cli/cmd/common"
	ident "github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/api/nchain"
	"github.com/provideservices/provide-go/api/vault"
	"github.com/spf13/cobra"
)

const defaultNChainBaselineNetworkID = "66d44f30-9092-4182-a3c4-bc02736d6ae5"
const defaultWorkgroupType = "baseline"

var name string
var networkID string

var vaultID string
var babyJubJubKeyID string
var secp256k1KeyID string
var hdwalletID string
var rsa4096Key string

var initBaselineWorkgroupCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize baseline workgroup",
	Long:  `Initialize and configure a new baseline workgroup`,
	Run:   initWorkgroup,
}

func authorizeApplicationContext() {
	common.AuthorizeApplicationContext()
	_, err := nchain.CreateWallet(common.ApplicationAccessToken, map[string]interface{}{
		"purpose": 44,
	})
	if err != nil {
		log.Printf("failed to initialize HD wallet; %s", err.Error())
		os.Exit(1)
	}
}

func initWorkgroup(cmd *cobra.Command, args []string) {
	common.NetworkID = networkID
	common.AuthorizeOrganizationContext()

	token := common.RequireUserAuthToken()
	application, err := ident.CreateApplication(token, map[string]interface{}{
		"config": map[string]interface{}{
			"baselined": true,
		},
		"name":       name,
		"network_id": common.NetworkID,
		"type":       defaultWorkgroupType,
	})
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	common.ApplicationID = application.ID.String()
	authorizeApplicationContext()

	common.InitWorkgroupContract()

	common.RequireOrganizationVault()
	requireOrganizationKeys()

	common.RegisterWorkgroupOrganization(application.ID.String())
	//common.RequireOrganizationEndpoints(nil)

	log.Printf("initialized baseline workgroup: %s", application.ID)
}

func requireOrganizationKeys() {
	var key *vault.Key
	var err error

	key, err = common.RequireOrganizationKeypair("babyJubJub")
	if err == nil {
		babyJubJubKeyID = key.ID.String()
	}

	key, err = common.RequireOrganizationKeypair("secp256k1")
	if err == nil {
		secp256k1KeyID = key.ID.String()
	}

	key, err = common.RequireOrganizationKeypair("BIP39")
	if err == nil {
		hdwalletID = key.ID.String()
	}

	key, err = common.RequireOrganizationKeypair("RSA-4096")
	if err == nil {
		rsa4096Key = key.ID.String()
	}
}

func init() {
	initBaselineWorkgroupCmd.Flags().StringVar(&name, "name", "", "name of the baseline workgroup")
	initBaselineWorkgroupCmd.MarkFlagRequired("name")

	initBaselineWorkgroupCmd.Flags().StringVar(&networkID, "network", defaultNChainBaselineNetworkID, "nchain network id of the baseline mainnet to use for this workgroup")

	initBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineWorkgroupCmd.MarkFlagRequired("organization")

	initBaselineWorkgroupCmd.Flags().StringVar(&common.MessagingEndpoint, "endpoint", "", "public messaging endpoint used for sending and receiving protocol messages")
}
