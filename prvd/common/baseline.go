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

package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	uuid "github.com/kthomas/go.uuid"
	"github.com/ory/viper"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/nchain"
	"github.com/provideplatform/provide-go/api/pgrok"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/provideplatform/provide-go/common"
	util "github.com/provideplatform/provide-go/common"
	commonutil "github.com/provideplatform/provide-go/common/util"
)

//                                         :os/`
//                                     ./ymNNNNNdo-
//                                  -odNNNNNNNNNNNNNy+.
//                                omNNNNNNNNNNNNNNNNNNNm.
//                                +dNNNNNNNNNNNNNNNNNNdo`
//                        :.      :..+yNNNNNNNNNNNNy+.-+.
//                        /oo/-`  sNms:`-ohNNNNms:`:smmy`
//                          .:+s+/.`:odNho-./+-.+yNdo-`-
//                              -/os+:`./ymmysdNh+..:+s+
//                                 `-/oo/-`-os:`-/oo/-`
//                                     .:+s+::+o+:.
//  `/:                                   `-//-`      :/` .:`
//  -NN                                               mN- sd:
//  -NN`:+o+/-      `-/+/:`-/-  -/+++-     `:/+/:`    mN- /+. ./--++/.      `:++/-
//  -NNmNhsymNd/  `sNNdyhmNmNy sNdssdNh` -hNmysymNy`  mN- dN: oNNNhymNh.  -hNmysymNs`
//  -NNs`    :NN/`dNy`    /NNy mNo..`//`-NN:    `oNm` mN- dN: oNd`   sNh :Nm:    `oNm`
//  -NN`      yNh:NN.      yNy -hmNNNms`oNNmmmmmmmmm: mN- dN: oNs    /Nh oNNmmmmmmmmm.
//  -NNh.   `+NN: dNy`    /NNy :/.``-NN-:NN+    .so-  mN- dN: oNs    /Nh :NN+`   .so.
//  -NNdNdhhmNh:  `smNdhhNNmNy +NNhhmNh` -yNNhhdNms`  mN- dN: oNs    /Nh  -hNNhhdNms`
//  `// .///:`      `-///-`:/:  `:///-     `:///-     //` //. -/-    ./:    `:///-`

// ██████╗  █████╗ ███████╗███████╗██╗     ██╗███╗   ██╗███████╗
// ██╔══██╗██╔══██╗██╔════╝██╔════╝██║     ██║████╗  ██║██╔════╝
// ██████╔╝███████║███████╗█████╗  ██║     ██║██╔██╗ ██║█████╗
// ██╔══██╗██╔══██║╚════██║██╔══╝  ██║     ██║██║╚██╗██║██╔══╝
// ██████╔╝██║  ██║███████║███████╗███████╗██║██║ ╚████║███████╗
// ╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═══╝╚══════╝

const defaultBaselineRegistryContractType = "registry"
const defaultBaselineOrgRegistryContractType = "organization-registry"

const defaultBaselineRegistryContractName = "Shuttle"
const defaultBaselineOrgRegistryContractName = "OrgRegistry"

const requireContractSleepInterval = time.Second * 1
const requireContractTickerInterval = time.Second * 5
const requireContractTimeout = time.Minute * 10

const requireOrganizationAPIEndpointTimeout = time.Second * 10
const requireOrganizationMessagingEndpointTimeout = time.Second * 10

var Tunnel bool
var APIEndpoint string
var ExposeAPITunnel bool
var ExposeMessagingTunnel bool
var ExposeWebsocketMessagingTunnel bool
var MessagingEndpoint string
var tunnelClient *pgrok.Client

var ApplicationAccessToken string
var OrganizationAccessToken string
var OrganizationRefreshToken string

var VaultID string

var ResolvedBaselineOrgAddress string // HACK

func AuthorizeApplicationContext() {
	RequireWorkgroup()

	token, err := ident.CreateToken(RequireUserAccessToken(), map[string]interface{}{
		"scope":          "offline_access",
		"application_id": ApplicationID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of application %s; %s", ApplicationID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		ApplicationAccessToken = *token.AccessToken
	}
}

func AuthorizeOrganizationContext(persist bool) {
	RequireOrganization()

	token, err := ident.CreateToken(RequireUserAccessToken(), map[string]interface{}{
		"scope":           "offline_access",
		"organization_id": OrganizationID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of organization %s; %s", OrganizationID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		OrganizationAccessToken = *token.AccessToken

		if token.RefreshToken != nil {
			OrganizationRefreshToken = *token.RefreshToken
		}

		if persist {
			// FIXME-- DRY this up (also exists in api_tokens_init.go)
			orgAPIAccessTokenKey := BuildConfigKeyWithID(AccessTokenConfigKey, OrganizationID)
			orgAPIRefreshTokenKey := BuildConfigKeyWithID(RefreshTokenConfigKey, OrganizationID)

			if token.AccessToken != nil {
				// fmt.Printf("Access token authorized for organization: %s\t%s\n", OrganizationID, *token.AccessToken)
				if !viper.IsSet(orgAPIAccessTokenKey) {
					viper.Set(orgAPIAccessTokenKey, *token.AccessToken)
					viper.WriteConfig()
				}
				if token.RefreshToken != nil {
					// fmt.Printf("Refresh token authorized for organization: %s\t%s\n", OrganizationID, *token.RefreshToken)
					if !viper.IsSet(orgAPIRefreshTokenKey) {
						viper.Set(orgAPIRefreshTokenKey, *token.RefreshToken)
						viper.WriteConfig()
					}
				}
			}
		}
	}
}

func InitWorkgroupContract(contractAddress string) *nchain.Contract {
	wallet, err := nchain.CreateWallet(OrganizationAccessToken, map[string]interface{}{
		"purpose": 44,
	})
	if err != nil {
		log.Printf("failed to initialize wallet for organization; %s", err.Error())
		os.Exit(1)
	}

	compiledArtifact := resolveBaselineOrgRegistryContractArtifact()
	contractName := defaultBaselineOrgRegistryContractName
	contractType := defaultBaselineOrgRegistryContractType

	if contractAddress == "0x" {
		compiledArtifact = resolveBaselineRegistryContractArtifact()
		contractName = defaultBaselineRegistryContractName
		contractType = defaultBaselineRegistryContractType
	}

	if compiledArtifact == nil {
		log.Printf("failed to resolve global baseline organization registry contract artifact")
		os.Exit(1)
	}

	log.Printf("deploying global baseline organization registry contract: %s", defaultBaselineRegistryContractName)
	contract, err := nchain.CreateContract(OrganizationAccessToken, map[string]interface{}{
		"address":    contractAddress,
		"name":       contractName,
		"network_id": NetworkID,
		"params": map[string]interface{}{
			"argv":              []interface{}{},
			"compiled_artifact": compiledArtifact,
			"wallet_id":         wallet.ID,
		},
		"type": contractType,
	})
	if err != nil {
		log.Printf("failed to initialize registry contract; %s", err.Error())
		os.Exit(1)
	}

	if contractAddress == "0x" {
		contract, err := RequireContract(nil, common.StringOrNil(defaultBaselineOrgRegistryContractType), true)
		if err != nil {
			log.Printf("failed to initialize registry contract; %s", err.Error())
			os.Exit(1)
		}

		return contract
	}

	return contract
}

func RegisterWorkgroupOrganization(workgroupID string) {
	_, err := RequireContract(nil, util.StringOrNil(defaultBaselineOrgRegistryContractType), false)
	if err != nil {
		log.Printf("failed to initialize registry contract; %s", err.Error())
		os.Exit(1)
	}
	err = ident.CreateApplicationOrganization(OrganizationAccessToken, workgroupID, map[string]interface{}{
		"organization_id": OrganizationID,
	})
	if err != nil {
		orgs, err := ident.ListApplicationOrganizations(OrganizationAccessToken, workgroupID, map[string]interface{}{
			"organization_id": OrganizationID,
		})
		if err == nil {
			// FIXME--
			for _, org := range orgs {
				if org.ID != nil && *org.ID == OrganizationID {
					return
				}
			}
			if len(orgs) > 0 && orgs[0].ID != nil && *orgs[0].ID == OrganizationID {
				return
			}
		}
		log.Printf("WARNING: organization not associated with workgroup")
		os.Exit(1)
	}
}

func RequireOrganizationVault() error {
	if OrganizationAccessToken == "" {
		return fmt.Errorf("organization access token not found")
	}

	// FIXME-- parameterize with --vault or similar?
	vaults, err := vault.ListVaults(OrganizationAccessToken, map[string]interface{}{})
	if err != nil {
		return err
	}

	if len(vaults) > 0 {
		VaultID = vaults[0].ID.String()
		return nil
	}

	vault, err := vault.CreateVault(OrganizationAccessToken, map[string]interface{}{
		"name":        fmt.Sprintf("vault for organization: %s", OrganizationID),
		"description": fmt.Sprintf("identity/signing keystore for organization: %s", OrganizationID),
	})
	if err != nil {
		return err
	}

	VaultID = vault.ID.String()
	return nil
}

func RequireOrganizationKeypair(spec string) (*vault.Key, error) {
	if VaultID == "" {
		RequireOrganizationVault()
	}

	// FIXME-- parameterize each key i.e. --secp256k1-key or similar?
	keys, err := vault.ListKeys(OrganizationAccessToken, VaultID, map[string]interface{}{
		"spec": spec,
	})
	if err != nil {
		log.Printf("failed to retrieve %s keys for organization: %s; %s", spec, OrganizationID, err.Error())
		return nil, err
	}

	if len(keys) > 0 {
		return keys[0], nil
	}

	key, err := vault.CreateKey(OrganizationAccessToken, VaultID, map[string]interface{}{
		"name":        fmt.Sprintf("%s key organization: %s", spec, OrganizationID),
		"description": fmt.Sprintf("%s key organization: %s", spec, OrganizationID),
		"spec":        spec,
		"type":        "asymmetric",
		"usage":       "sign/verify",
	})
	if err != nil {
		return nil, err
	}

	return key, nil
}

func RequireContract(contractID, contractType *string, printCreationTxLink bool) (*nchain.Contract, error) {
	startTime := time.Now()
	timer := time.NewTicker(requireContractTickerInterval)

	printed := false

	for {
		select {
		case <-timer.C:
			var contract *nchain.Contract
			var err error
			if contractID != nil {
				contract, err = nchain.GetContractDetails(OrganizationAccessToken, *contractID, map[string]interface{}{})
			} else if contractType != nil {
				contracts, _ := nchain.ListContracts(OrganizationAccessToken, map[string]interface{}{
					"type": contractType,
				})
				if len(contracts) > 0 {
					contract = contracts[0]
				}
			}

			// FIXME-- KT-- review removal of contract.TransactionID != nil condition
			if err == nil && contract != nil {
				if !printed && printCreationTxLink && contract.TransactionID != nil {
					tx, _ := nchain.GetTransactionDetails(OrganizationAccessToken, contract.TransactionID.String(), map[string]interface{}{})
					if tx.Hash != nil {
						etherscanBaseURL := EtherscanBaseURL(tx.NetworkID.String())
						if etherscanBaseURL != nil {
							log.Printf("View on Etherscan: %s/tx/%s", *etherscanBaseURL, *tx.Hash) // HACK
						} else {
							log.Printf("Transaction hash: %s", *tx.Hash)
						}
						printed = true
					}
				}

				if contract.Address != nil && *contract.Address != "0x" {
					if Verbose {
						tx, _ := nchain.GetTransactionDetails(OrganizationAccessToken, contract.TransactionID.String(), map[string]interface{}{})
						txraw, _ := json.MarshalIndent(tx, "", "  ")
						log.Printf("%s", string(txraw))
					}

					return contract, nil
				}
			}
		default:
			if startTime.Add(requireContractTimeout).Before(time.Now()) {
				log.Printf("WARNING: workgroup contract deployment timed out")
				return nil, errors.New("workgroup contract deployment timed out")
			} else {
				time.Sleep(requireContractSleepInterval)
			}
		}
	}
}

func resolveBaselineRegistryContractArtifact() *nchain.CompiledArtifact {
	capabilities, err := commonutil.ResolveCapabilitiesManifest()
	if err != nil {
		return nil
	}

	var registryArtifact *nchain.CompiledArtifact
	if baseline, baselineOk := capabilities["baseline"].(map[string]interface{}); baselineOk {
		if contracts, contractsOk := baseline["contracts"].([]interface{}); contractsOk {
			for _, contract := range contracts {
				isShuttleContract := false
				if name, nameOk := contract.(map[string]interface{})["name"].(string); nameOk && strings.ToLower(name) == "shuttle" {
					isShuttleContract = true
				} else if name, nameOk := contract.(map[string]interface{})["contractName"].(string); nameOk && strings.ToLower(name) == "shuttle" {
					isShuttleContract = true
				}

				if isShuttleContract {
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

func resolveBaselineOrgRegistryContractArtifact() *nchain.CompiledArtifact {
	capabilities, err := commonutil.ResolveCapabilitiesManifest()
	if err != nil {
		return nil
	}

	var registryArtifact *nchain.CompiledArtifact
	if baseline, baselineOk := capabilities["baseline"].(map[string]interface{}); baselineOk {
		if contracts, contractsOk := baseline["contracts"].([]interface{}); contractsOk {
			var artifact nchain.CompiledArtifact
			raw, _ := json.Marshal(contracts[1])
			json.Unmarshal(raw, &artifact)

			return &artifact
		}
	}

	return registryArtifact
}

// RequireOrganizationEndpoints fn is the function to call after the tunnel has been established,
// prior to the runloop and signal handling is installed
func RequireOrganizationEndpoints(fn func(), tunnelShutdownFn func(*string), apiPort, messagingPort, websocketMessagingPort int) {
	run := func() {
		if OrganizationID == "" {
			fmt.Println("WARNING: failed to set organization endpoints; organization id not set")
			os.Exit(1)
		}

		if WorkgroupID == "" {
			fmt.Println("WARNING: failed to set organization endpoints; workgroup id not set")
			os.Exit(1)
		}

		wgID, err := uuid.FromString(WorkgroupID)
		if err != nil {
			log.Printf("WARNING: failed to update organization; %s", err.Error())
			os.Exit(1)
		}

		if Organization.Metadata == nil {
			Organization.Metadata = &OrganizationMetadata{}
		}

		if Organization.Metadata.Workgroups == nil {
			Organization.Metadata.Workgroups = map[uuid.UUID]*OrganizationWorkgroupMetadata{}
		}

		if Organization.Metadata.Workgroups[wgID] == nil {
			Organization.Metadata.Workgroups[wgID] = &OrganizationWorkgroupMetadata{}
		}

		key, err := RequireOrganizationKeypair("secp256k1")
		if err != nil {
			log.Printf("WARNING: failed to update organization; %s", err.Error())
			os.Exit(1)
		}

		Organization.Metadata.Address = *key.Address
		ResolvedBaselineOrgAddress = *key.Address

		if APIEndpoint != "" {
			Organization.Metadata.BPIEndpoint = APIEndpoint
		} else {
			Organization.Metadata.BPIEndpoint = "http://localhost:8080"
		}

		if MessagingEndpoint != "" {
			Organization.Metadata.MessagingEndpoint = MessagingEndpoint
		} else {
			Organization.Metadata.MessagingEndpoint = "nats://localhost:4222"
		}

		Organization.Metadata.Domain = "baseline.local"

		Organization.Metadata.Workgroups[wgID].BPIEndpoint = &APIEndpoint

		var org map[string]interface{}
		raw, _ := json.Marshal(Organization)
		json.Unmarshal(raw, &org)

		if err := ident.UpdateOrganization(OrganizationAccessToken, OrganizationID, org); err != nil {
			log.Printf("WARNING: failed to update organization; %s", err.Error())
			os.Exit(1)
		}

		log.Printf("successfully set BPI endpoint %s on organization %s\n", APIEndpoint, OrganizationID)

		if fn != nil {
			fn()
		}
	}

	if Tunnel {
		ExposeAPITunnel = true
		ExposeMessagingTunnel = true
	}

	if !ExposeAPITunnel && !ExposeMessagingTunnel {
		publicIP, err := util.ResolvePublicIP()
		if err != nil {
			log.Printf("WARNING: failed to resolve public IP")
			os.Exit(1)
		}

		APIEndpoint = fmt.Sprintf("http://%s:%d", *publicIP, apiPort)
		MessagingEndpoint = fmt.Sprintf("nats://%s:%d", *publicIP, messagingPort)

		run()
	} else {
		const runloopSleepInterval = 250 * time.Millisecond
		const runloopTickInterval = 5000 * time.Millisecond

		var (
			cancelF     context.CancelFunc
			closing     uint32
			shutdownCtx context.Context
			sigs        chan os.Signal
		)

		installSignalHandlers := func() {
			log.Printf("installing signal handlers")
			sigs = make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			shutdownCtx, cancelF = context.WithCancel(context.Background())
		}

		shutdown := func() {
			if atomic.AddUint32(&closing, 1) == 1 {
				log.Print("shutting down")
				tunnelClient.Close()
				cancelF()
				os.Exit(0)
			}
		}

		shuttingDown := func() bool {
			return (atomic.LoadUint32(&closing) > 0)
		}

		installSignalHandlers()

		var once sync.Once
		_tunnelShutdownFn := func(reason *string) {
			once.Do(func() {
				if tunnelShutdownFn != nil {
					tunnelShutdownFn(reason)
				}
				shutdown()
			})
		}

		go func() {
			var err error
			tunnelClient, err = pgrok.Factory()
			if err != nil {
				log.Printf("WARNING: failed to initialize tunnel; %s", err.Error())
				os.Exit(1)
			}

			if ExposeAPITunnel {
				tunnel, _ := tunnelClient.TunnelFactory(
					fmt.Sprintf("%s-api", OrganizationID),
					fmt.Sprintf("127.0.0.1:%d", apiPort),
					nil,
					common.StringOrNil("https"),
					common.StringOrNil(OrganizationAccessToken),
					_tunnelShutdownFn,
				)
				tunnelClient.AddTunnel(tunnel)
			}

			if ExposeMessagingTunnel {
				tunnel, _ := tunnelClient.TunnelFactory(
					fmt.Sprintf("%s-msg", OrganizationID),
					fmt.Sprintf("127.0.0.1:%d", messagingPort),
					nil,
					common.StringOrNil("tcp"),
					common.StringOrNil(OrganizationAccessToken),
					_tunnelShutdownFn,
				)
				tunnelClient.AddTunnel(tunnel)
			}

			if ExposeWebsocketMessagingTunnel {
				tunnel, _ := tunnelClient.TunnelFactory(
					fmt.Sprintf("%s-wss", OrganizationID),
					fmt.Sprintf("127.0.0.1:%d", websocketMessagingPort),
					nil,
					common.StringOrNil("https"),
					common.StringOrNil(OrganizationAccessToken),
					_tunnelShutdownFn,
				)
				tunnelClient.AddTunnel(tunnel)
			}

			err = tunnelClient.ConnectAll()
			if err != nil {
				log.Printf("WARNING: failed to initialize tunnel(s); %s", err.Error())
				os.Exit(1)
			}

			if ExposeAPITunnel {
				for tunnelClient.Tunnels[0].RemoteAddr == nil {
					time.Sleep(time.Millisecond * 10)
				}

				APIEndpoint = *tunnelClient.Tunnels[0].RemoteAddr
				log.Printf("established tunnel connection for API endpoint: %s\n", APIEndpoint)
			}

			if ExposeMessagingTunnel {
				i := len(tunnelClient.Tunnels) - 1
				if ExposeWebsocketMessagingTunnel {
					i--
				}

				for tunnelClient.Tunnels[i].RemoteAddr == nil {
					time.Sleep(time.Millisecond * 10)
				}

				MessagingEndpoint = *tunnelClient.Tunnels[i].RemoteAddr
				log.Printf("established tunnel connection for messaging endpoint: %s\n", MessagingEndpoint)
			}
		}()

		if ExposeAPITunnel {
			go func() {
				startTime := time.Now()
				for APIEndpoint == "" {
					if startTime.Add(requireOrganizationAPIEndpointTimeout).Before(time.Now()) {
						log.Printf("WARNING: organization API endpoint tunnel timed out")
						os.Exit(1)
					}
					time.Sleep(time.Millisecond * 10)
				}
			}()
		}

		if ExposeMessagingTunnel {
			go func() {
				startTime := time.Now()
				for MessagingEndpoint == "" {
					if startTime.Add(requireOrganizationMessagingEndpointTimeout).Before(time.Now()) {
						log.Printf("WARNING: organization messaging endpoint tunnel timed out")
						os.Exit(1)
					}
					time.Sleep(time.Millisecond * 10)
				}
			}()
		}

		for (ExposeAPITunnel && APIEndpoint == "") || (ExposeMessagingTunnel && MessagingEndpoint == "") {
			time.Sleep(time.Millisecond * 10)
		}

		run()

		// log.Printf("starting tunnel runloop")
		timer := time.NewTicker(runloopTickInterval)
		defer timer.Stop()

		for !shuttingDown() {
			select {
			case <-timer.C:
				// tick... no-op
			case sig := <-sigs:
				fmt.Printf("received signal: %s", sig)
				shutdown()
			case <-shutdownCtx.Done():
				close(sigs)
			default:
				time.Sleep(runloopSleepInterval)
			}
		}

		log.Printf("exiting tunnel runloop")
		cancelF()
	}
}
