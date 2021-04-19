package common

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kthomas/gonnel"
	"github.com/provideservices/provide-go/api"
	"github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/api/nchain"
	"github.com/provideservices/provide-go/api/vault"
	util "github.com/provideservices/provide-go/common"
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

const defaultBaselineRegistryContractName = "Shuttle"

const requireContractSleepInterval = time.Second * 1
const requireContractTickerInterval = time.Second * 5
const requireContractTimeout = time.Minute * 10

const requireOrganizationMessagingEndpointTimeout = time.Second * 5

var ExposeTunnel bool
var MessagingEndpoint string
var tunnelClient *gonnel.Client

var ApplicationAccessToken string
var OrganizationAccessToken string

var VaultID string

func AuthorizeApplicationContext() {
	token, err := ident.CreateToken(RequireUserAuthToken(), map[string]interface{}{
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

func AuthorizeOrganizationContext() {
	token, err := ident.CreateToken(RequireUserAuthToken(), map[string]interface{}{
		"scope":           "offline_access",
		"organization_id": OrganizationID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of organization %s; %s", OrganizationID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		OrganizationAccessToken = *token.AccessToken
	}
}

func InitWorkgroupContract() *nchain.Contract {
	wallet, err := nchain.CreateWallet(OrganizationAccessToken, map[string]interface{}{
		"purpose": 44,
	})
	if err != nil {
		log.Printf("failed to initialize wallet for organization; %s", err.Error())
		os.Exit(1)
	}

	log.Printf("deploying global baseline organization registry contract: %s", defaultBaselineRegistryContractName)
	contract, err := nchain.CreateContract(ApplicationAccessToken, map[string]interface{}{
		"address":    "0x",
		"name":       defaultBaselineRegistryContractName,
		"network_id": NetworkID,
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

	err = RequireContract(util.StringOrNil(contract.ID.String()), nil)
	if err != nil {
		log.Printf("failed to initialize registry contract; %s", err.Error())
		os.Exit(1)
	}

	return contract
}

func RegisterWorkgroupOrganization(applicationID string) {
	err := RequireContract(nil, util.StringOrNil("organization-registry"))
	if err != nil {
		log.Printf("failed to initialize registry contract; %s", err.Error())
		os.Exit(1)
	}
	err = ident.CreateApplicationOrganization(ApplicationAccessToken, applicationID, map[string]interface{}{
		"organization_id": OrganizationID,
	})
	if err != nil {
		orgs, err := ident.ListApplicationOrganizations(ApplicationAccessToken, applicationID, map[string]interface{}{
			"organization_id": OrganizationID,
		})
		if err == nil {
			// FIXME--
			for _, org := range orgs {
				if org.ID.String() == OrganizationID {
					return
				}
			}
			if len(orgs) > 0 && orgs[0].ID.String() == OrganizationID {
				return
			}
		} else {
			log.Printf("WARNING: organization not associated with workgroup")
			os.Exit(1)
		}
	}
}

// fn is the function to call after the tunnel has been established,
// prior to the runloop and signal handling is installed
func RequireOrganizationMessagingEndpoint(fn func()) {
	setupMessagingEndpoint(fn)
}

func RequireOrganizationVault() {
	// FIXME-- parameterize with --vault or similar?
	vaults, err := vault.ListVaults(OrganizationAccessToken, map[string]interface{}{
		"organization_id": OrganizationID,
	})
	if err != nil {
		log.Printf("failed to retrieve vaults for organization: %s; %s", OrganizationID, err.Error())
		os.Exit(1)
	}

	if len(vaults) > 0 {
		VaultID = vaults[0].ID.String()
	} else {
		vault, err := vault.CreateVault(OrganizationAccessToken, map[string]interface{}{
			"name":        fmt.Sprintf("vault for organization: %s", OrganizationID),
			"description": fmt.Sprintf("identity/signing keystore for organization: %s", OrganizationID),
		})
		if err != nil {
			log.Printf("failed to create vault for organization: %s; %s", OrganizationID, err.Error())
			os.Exit(1)
		}
		VaultID = vault.ID.String()
	}

	if VaultID == "" {
		log.Printf("failed to require vault for organization: %s", OrganizationID)
		os.Exit(1)
	}
}

func RequireOrganizationKeypair(spec string) (*vault.Key, error) {
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

func ResolveCapabilities() (map[string]interface{}, error) {
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

func RequireContract(contractID, contractType *string) error {
	startTime := time.Now()
	timer := time.NewTicker(requireContractTickerInterval)

	for {
		select {
		case <-timer.C:
			var contract *nchain.Contract
			var err error
			if contractID != nil {
				contract, err = nchain.GetContractDetails(ApplicationAccessToken, *contractID, map[string]interface{}{})
			} else if contractType != nil {
				contracts, _ := nchain.ListContracts(ApplicationAccessToken, map[string]interface{}{
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

func resolveBaselineRegistryContractArtifact() *nchain.CompiledArtifact {
	capabilities, err := ResolveCapabilities()
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

// fn is the function to call after the tunnel has been established,
// prior to the runloop and signal handling is installed
func setupMessagingEndpoint(fn func()) {
	if MessagingEndpoint == "" && !ExposeTunnel {
		publicIP, err := util.ResolvePublicIP()
		if err != nil {
			log.Printf("WARNING: failed to resolve public IP")
			os.Exit(1)
		}

		MessagingEndpoint = fmt.Sprintf("nats://%s:4222", *publicIP)
	} else if ExposeTunnel {
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
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			shutdownCtx, cancelF = context.WithCancel(context.Background())
		}

		shutdown := func() {
			if atomic.AddUint32(&closing, 1) == 1 {
				log.Print("shutting down")
				cancelF()
				os.Exit(0)
			}
		}

		shuttingDown := func() bool {
			return (atomic.LoadUint32(&closing) > 0)
		}

		installSignalHandlers()

		go func() {
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

				defer tunnelClient.Close()

				done := make(chan bool)
				go tunnelClient.StartServer(done)
				<-done

				tunnelClient.AddTunnel(&gonnel.Tunnel{
					Proto:        gonnel.TCP,
					Name:         fmt.Sprintf("%s-endpoint", OrganizationID),
					LocalAddress: "127.0.0.1:4222", // FIXME-- this port
				})

				tunnelClient.ConnectAll()

				MessagingEndpoint = tunnelClient.Tunnel[0].RemoteAddress
				log.Printf("established tunnel connection for messaging endpoint: %s\n", MessagingEndpoint)

				fmt.Print("Press any key to disconnect...\n")
				reader := bufio.NewReader(os.Stdin)
				reader.ReadRune()

				tunnelClient.DisconnectAll()
			}
		}()

		if fn != nil {
			go fn()
		}

		if ExposeTunnel {
			startTime := time.Now()
			for MessagingEndpoint == "" {
				time.Sleep(time.Millisecond * 50)
				if startTime.Add(requireOrganizationMessagingEndpointTimeout).Before(time.Now()) {
					log.Printf("WARNING: organization messaging endpoint tunnel timed out")
					os.Exit(1)
				}
			}
		}

		org, err := ident.GetOrganizationDetails(OrganizationAccessToken, OrganizationID, map[string]interface{}{})
		if err != nil {
			log.Printf("failed to retrieve organization: %s; %s", OrganizationID, err.Error())
			os.Exit(1)
		}

		if org.Metadata == nil {
			org.Metadata = map[string]interface{}{}
		}
		org.Metadata["messaging_endpoint"] = MessagingEndpoint

		err = ident.UpdateOrganization(RequireUserAuthToken(), OrganizationID, map[string]interface{}{
			"metadata": org.Metadata,
		})
		if err != nil {
			log.Printf("failed to update messaging endpoint for organization: %s; %s", OrganizationID, err.Error())
			os.Exit(1)
		}

		log.Printf("starting tunnel runloop")

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
