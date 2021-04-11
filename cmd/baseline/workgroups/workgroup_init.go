package workgroups

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/provideservices/provide-cli/cmd/common"
	api "github.com/provideservices/provide-go/api"
	ident "github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/api/nchain"
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

	token := common.RequireAPIToken()
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

	return nil
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

func init() {
	initBaselineWorkgroupCmd.Flags().StringVar(&name, "name", "", "name of the baseline workgroup")
	initBaselineWorkgroupCmd.MarkFlagRequired("name")

	initBaselineWorkgroupCmd.Flags().StringVar(&networkID, "network", defaultNChainBaselineNetworkID, "nchain network id of the baseline mainnet to use for this workgroup")

	initBaselineWorkgroupCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineWorkgroupCmd.MarkFlagRequired("organization")
}
