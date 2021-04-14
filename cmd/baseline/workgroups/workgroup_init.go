package workgroups

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kthomas/gonnel"
	"github.com/provideservices/provide-cli/cmd/common"
	api "github.com/provideservices/provide-go/api"
	ident "github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/api/nchain"
	"github.com/provideservices/provide-go/api/vault"
	util "github.com/provideservices/provide-go/common"
	"github.com/spf13/cobra"
)

const defaultNChainBaselineNetworkID = "66d44f30-9092-4182-a3c4-bc02736d6ae5"
const defaultWorkgroupType = "baseline"

const defaultBaselineRegistryContractName = "Shuttle"

const requireContractSleepInterval = time.Second * 1
const requireContractTickerInterval = time.Second * 5
const requireContractTimeout = time.Minute * 10

var name string
var networkID string
var organizationAccessToken string
var applicationAccessToken string

var messagingEndpoint string
var exposeTunnel bool
var tunnelClient *gonnel.Client

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
	token, err := ident.CreateToken(common.RequireUserAuthToken(), map[string]interface{}{
		"scope":          "offline_access",
		"application_id": common.ApplicationID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of application %s; %s", common.ApplicationID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		applicationAccessToken = *token.AccessToken
	}

	// HACK...
	_, err = nchain.CreateWallet(applicationAccessToken, map[string]interface{}{
		"purpose": 44,
	})
}

func authorizeOrganizationContext() {
	token, err := ident.CreateToken(common.RequireUserAuthToken(), map[string]interface{}{
		"scope":           "offline_access",
		"organization_id": common.OrganizationID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of organization %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		organizationAccessToken = *token.AccessToken
	}
}

func initWorkgroup(cmd *cobra.Command, args []string) {
	authorizeOrganizationContext()

	token := common.RequireUserAuthToken()
	application, err := ident.CreateApplication(token, map[string]interface{}{
		"config": map[string]interface{}{
			"baselined": true,
		},
		"name":       name,
		"network_id": networkID,
		"type":       defaultWorkgroupType,
	})
	if err != nil {
		log.Printf("failed to initialize baseline workgroup; %s", err.Error())
		os.Exit(1)
	}

	common.ApplicationID = application.ID.String()
	authorizeApplicationContext()

	initWorkgroupContract()

	requireOrganizationVault()
	requireOrganizationKeys()
	requireOrganizationMessagingEndpoint()
	registerWorkgroupOrganization(application.ID.String())

	log.Printf("initialized baseline workgroup: %s", application.ID)
}

func initWorkgroupContract() *nchain.Contract {
	wallet, err := nchain.CreateWallet(organizationAccessToken, map[string]interface{}{
		"purpose": 44,
	})
	if err != nil {
		log.Printf("failed to initialize wallet for organization; %s", err.Error())
		os.Exit(1)
	}

	log.Printf("deploying global baseline organization registry contract: %s", defaultBaselineRegistryContractName)
	contract, err := nchain.CreateContract(applicationAccessToken, map[string]interface{}{
		"address":    "0x",
		"name":       defaultBaselineRegistryContractName,
		"network_id": networkID,
		"params": map[string]interface{}{
			"argv":              []interface{}{},
			"compiled_artifact": resolveBaselineRegistryContractArtifact(),
			"wallet_id":         wallet.ID,
		},
		"type": "registry",
	})
	if err != nil {
		log.Printf("failed to initialize registry contract; %s", err.Error())
		os.Exit(1)
	}

	err = requireContract(util.StringOrNil(contract.ID.String()), nil)
	if err != nil {
		log.Printf("failed to initialize registry contract; %s", err.Error())
		os.Exit(1)
	}

	return contract
}

func registerWorkgroupOrganization(applicationID string) {
	err := requireContract(nil, util.StringOrNil("organization-registry"))
	if err != nil {
		log.Printf("failed to initialize registry contract; %s", err.Error())
		os.Exit(1)
	}
	err = ident.CreateApplicationOrganization(applicationAccessToken, applicationID, map[string]interface{}{
		"organization_id": common.OrganizationID,
	})
	if err != nil {
		log.Printf("WARNING: organization not associated with workgroup")
		os.Exit(1)
	}
}

func requireContract(contractID, contractType *string) error {
	startTime := time.Now()
	timer := time.NewTicker(requireContractTickerInterval)

	for {
		select {
		case <-timer.C:
			var contract *nchain.Contract
			var err error
			if contractID != nil {
				contract, err = nchain.GetContractDetails(applicationAccessToken, *contractID, map[string]interface{}{})
			} else if contractType != nil {
				contracts, _ := nchain.ListContracts(applicationAccessToken, map[string]interface{}{
					"type": contractType,
				})
				if len(contracts) > 0 {
					contract = contracts[0]
				}
			}

			if err == nil && contract != nil {
				if contract.Address != nil && *contract.Address != "0x" {
					return nil
				}
			}
		default:
			if startTime.Add(requireContractTimeout).Before(time.Now()) {
				log.Printf("WARNING: workgroup contract deployment timed out")
				return errors.New("workgroup contract deployment timed out")
			} else {
				time.Sleep(requireContractSleepInterval)
			}
		}
	}
}

func requireOrganizationMessagingEndpoint() {
	setupMessagingEndpoint()
	if exposeTunnel {
		for messagingEndpoint == "" {
			time.Sleep(time.Millisecond * 50)
		}
	}

	org, err := ident.GetOrganizationDetails(organizationAccessToken, common.OrganizationID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to retrieve organization: %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}

	if org.Metadata == nil {
		org.Metadata = map[string]interface{}{
			"messaging_endpoint": messagingEndpoint,
		}
	}

	err = ident.UpdateOrganization(common.RequireUserAuthToken(), common.OrganizationID, map[string]interface{}{
		"metadata": org.Metadata,
	})
	if err != nil {
		log.Printf("failed to update messaging endpoint for organization: %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}
	log.Printf("messaging endpoint set: %s", messagingEndpoint)
}

func requireOrganizationVault() {
	// FIXME-- parameterize with --vault or similar?
	vaults, err := vault.ListVaults(organizationAccessToken, map[string]interface{}{
		"organization_id": common.OrganizationID,
	})
	if err != nil {
		log.Printf("failed to retrieve vaults for organization: %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}

	if len(vaults) > 0 {
		vaultID = vaults[0].ID.String()
	} else {
		vault, err := vault.CreateVault(organizationAccessToken, map[string]interface{}{
			"name":        fmt.Sprintf("vault for organization: %s", common.OrganizationID),
			"description": fmt.Sprintf("identity/signing keystore for organization: %s", common.OrganizationID),
		})
		if err != nil {
			log.Printf("failed to create vault for organization: %s; %s", common.OrganizationID, err.Error())
			os.Exit(1)
		}
		vaultID = vault.ID.String()
	}

	if vaultID == "" {
		log.Printf("failed to require vault for organization: %s", common.OrganizationID)
		os.Exit(1)
	}
}

func requireOrganizationKeys() {
	var key *vault.Key
	var err error

	key, err = requireOrganizationKeypair("babyJubJub")
	if err == nil {
		babyJubJubKeyID = key.ID.String()
	}

	key, err = requireOrganizationKeypair("secp256k1")
	if err == nil {
		secp256k1KeyID = key.ID.String()
	}

	key, err = requireOrganizationKeypair("BIP39")
	if err == nil {
		hdwalletID = key.ID.String()
	}

	key, err = requireOrganizationKeypair("RSA-4096")
	if err == nil {
		rsa4096Key = key.ID.String()
	}
}

func requireOrganizationKeypair(spec string) (*vault.Key, error) {
	// FIXME-- parameterize each key i.e. --secp256k1-key or similar?
	keys, err := vault.ListKeys(organizationAccessToken, vaultID, map[string]interface{}{
		"spec": spec,
	})
	if err != nil {
		log.Printf("failed to retrieve %s keys for organization: %s; %s", spec, common.OrganizationID, err.Error())
		return nil, err
	}

	if len(keys) > 0 {
		return keys[0], nil
	}

	key, err := vault.CreateKey(organizationAccessToken, vaultID, map[string]interface{}{
		"name":        fmt.Sprintf("%s key organization: %s", spec, common.OrganizationID),
		"description": fmt.Sprintf("%s key organization: %s", spec, common.OrganizationID),
		"spec":        spec,
		"type":        "asymmetric",
		"usage":       "sign/verify",
	})
	if err != nil {
		return nil, err
	}

	return key, nil
}

func resolveCapabilities() (map[string]interface{}, error) {
	capabilitiesClient := &api.Client{
		Host:   "s3.amazonaws.com",
		Scheme: "https",
		Path:   "static.provide.services/capabilities",
	}
	_, capabilities, err := capabilitiesClient.Get("provide-capabilities-manifest.json", map[string]interface{}{})
	if err != nil {
		log.Printf("WARNING: failed to fetch capabilities; %s", err.Error())
		return nil, err
	}

	return capabilities.(map[string]interface{}), nil
}

func resolveBaselineRegistryContractArtifact() *nchain.CompiledArtifact {
	capabilities, err := resolveCapabilities()
	if err != nil {
		return nil
	}

	var registryArtifact *nchain.CompiledArtifact
	if baseline, baselineOk := capabilities["baseline"].(map[string]interface{}); baselineOk {
		if contracts, contractsOk := baseline["contracts"].([]interface{}); contractsOk {
			for _, contract := range contracts {
				if name, nameOk := contract.(map[string]interface{})["name"].(string); nameOk && strings.ToLower(name) == "shuttle" {
					raw, _ := json.Marshal(contract)
					err := json.Unmarshal(raw, &registryArtifact)
					if err != nil {
						log.Printf("failed to parse registry contract from capabilities; %s", err.Error())
						return nil
					}
				}
			}
		}
	}

	return registryArtifact
}

func setupMessagingEndpoint() {
	if messagingEndpoint == "" && !exposeTunnel {
		publicIP, err := util.ResolvePublicIP()
		if err != nil {
			log.Printf("WARNING: failed to resolve public IP")
			os.Exit(1)
		}

		messagingEndpoint = fmt.Sprintf("nats://%s:4222", *publicIP)
	} else if exposeTunnel {
		var out bytes.Buffer
		cmd := exec.Command("which", "ngrok")
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			tunnelClient, err = gonnel.NewClient(gonnel.Options{
				BinaryPath: strings.Trim(out.String(), "\n"),
			})
			if err != nil {
				log.Printf("WARNING: failed to initialize tunnel; %s", err.Error())
				os.Exit(1)
			}

			done := make(chan bool)
			go tunnelClient.StartServer(done)
			<-done

			tunnelClient.AddTunnel(&gonnel.Tunnel{
				Proto:        gonnel.TCP,
				Name:         fmt.Sprintf("%s-endpoint", common.OrganizationID),
				LocalAddress: "127.0.0.1:4222", // FIXME-- this port
			})

			tunnelClient.ConnectAll()

			messagingEndpoint = tunnelClient.Tunnel[0].RemoteAddress
			log.Printf("established tunnel connection for messaging endpoint: %s", messagingEndpoint)

			// fmt.Print("Press any to disconnect")
			// reader := bufio.NewReader(os.Stdin)
			// reader.ReadRune()

			// client.DisconnectAll()
		}
	}
}

func init() {
	initBaselineWorkgroupCmd.Flags().StringVar(&name, "name", "", "name of the baseline workgroup")
	initBaselineWorkgroupCmd.MarkFlagRequired("name")

	initBaselineWorkgroupCmd.Flags().StringVar(&networkID, "network", defaultNChainBaselineNetworkID, "nchain network id of the baseline mainnet to use for this workgroup")

	initBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineWorkgroupCmd.MarkFlagRequired("organization")

	initBaselineWorkgroupCmd.Flags().StringVar(&messagingEndpoint, "endpoint", "", "public messaging endpoint used for sending and receiving protocol messages")
	initBaselineWorkgroupCmd.Flags().BoolVar(&exposeTunnel, "tunnel", false, "when true, a tunnel established to expose this endpoint to the WAN")
}
