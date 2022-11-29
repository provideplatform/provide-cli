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

package stack

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/nchain"
	"github.com/provideplatform/provide-go/api/vault"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const baselineContainerImage = "provide/baseline"
const localBaselineContainerImage = "baseline_dev"
const identContainerImage = "provide/ident"
const nchainContainerImage = "provide/nchain"
const privacyContainerImage = "provide/privacy"
const vaultContainerImage = "provide/vault"
const elasticContainerImage = "docker.elastic.co/elasticsearch/elasticsearch:7.17.7"
const postgresContainerImage = "postgres"
const natsContainerImage = "provide/nats-server:2.7.2-PRVD"
const redisContainerImage = "redis"
const defaultContainerReachabilityTimeout = time.Millisecond * 2500
const defaultBPIStackName = "baseline-local"
const defaultElasticReachabilityTimeout = time.Second * 5
const defaultNatsServerName = "prvd"
const defaultNatsReachabilityTimeout = time.Second * 5
const defaultPostgresReachabilityTimeout = time.Second * 5
const defaultRedisReachabilityTimeout = time.Second * 5

const defaultJWTSignerPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAullT/WoZnxecxKwQFlwE
9lpQrekSD+txCgtb9T3JvvX/YkZTYkerf0rssQtrwkBlDQtm2cB5mHlRt4lRDKQy
EA2qNJGM1Yu379abVObQ9ZXI2q7jTBZzL/Yl9AgUKlDIAXYFVfJ8XWVTi0l32Vsx
tJSd97hiRXO+RqQu5UEr3jJ5tL73iNLp5BitRBwa4KbDCbicWKfSH5hK5DM75EyM
R/SzR3oCLPFNLs+fyc7zH98S1atglbelkZsMk/mSIKJJl1fZFVCUxA+8CaPiKbpD
QLpzydqyrk/y275aSU/tFHidoewvtWorNyFWRnefoWOsJFlfq1crgMu2YHTMBVtU
SJ+4MS5D9fuk0queOqsVUgT7BVRSFHgDH7IpBZ8s9WRrpE6XOE+feTUyyWMjkVgn
gLm5RSbHpB8Wt/Wssy3VMPV3T5uojPvX+ITmf1utz0y41gU+iZ/YFKeNN8WysLxX
AP3Bbgo+zNLfpcrH1Y27WGBWPtHtzqiafhdfX6LQ3/zXXlNuruagjUohXaMltH+S
K8zK4j7n+BYl+7y1dzOQw4CadsDi5whgNcg2QUxuTlW+TQ5VBvdUl9wpTSygD88H
xH2b0OBcVjYsgRnQ9OZpQ+kIPaFhaWChnfEArCmhrOEgOnhfkr6YGDHFenfT3/RA
PUl1cxrvY7BHh4obNa6Bf8ECAwEAAQ==
-----END PUBLIC KEY-----`

const defaultNATSStreamingClusterID = "provide"

const dockerNetworkDriverBridge = "bridge"
const dockerNetworkDriverDefault = "default"
const dockerNetworkDriverNAT = "nat"

const defaultBPIContainerPort = 8080
const defaultElasticContainerPort = 9200
const defaultNATSContainerPort = 4222
const defaultNATSWebsocketContainerPort = 4221
const defaultPostgresContainerPort = 5432
const defaultRedisContainerPort = 6379

type portMapping struct {
	hostPort      int
	containerPort int
}

var dockerNetworkID string
var Optional bool
var name string
var containerPort int
var port int
var identPort int
var nchainPort int
var privacyPort int
var vaultPort int

var elasticHostname string
var elasticContainerPort int
var elasticPort int
var elasticUsername string
var elasticPassword string
var elasticAPIScheme string
var elasticAcceptSelfSignedCertificate bool
var elasticMemory string

var natsContainerPort int
var natsPort int
var natsWebsocketContainerPort int
var natsWebsocketPort int
var natsWebsocketTLS bool
var postgresContainerPort int
var postgresPort int
var redisContainerPort int
var redisPort int

var apiHostname string
var consumerHostname string
var identHostname string
var identConsumerHostname string
var nchainHostname string
var nchainConsumerHostname string
var nchainStatsdaemonHostname string
var nchainReachabilitydaemonHostname string
var privacyHostname string
var privacyConsumerHostname string
var vaultHostname string
var natsHostname string
var natsServerName string
var postgresDatabase string
var postgresHostname string
var postgresPassword string
var postgresUser string
var redisHostname string
var redisHosts string

var autoRemove bool
var prune bool

var databaseLogging string
var logLevel string
var syslogEndpoint string

var baselineOrganizationAddress string

// var baselineOrganizationAPIEndpoint string
var baselineRegistryContractAddress string

var nchainBaselineNetworkID string

var jwtSignerPublicKey string
var natsAuthToken string

var identAPIHost string
var identAPIScheme string
var identContainerPort int

var nchainAPIHost string
var nchainAPIScheme string
var nchainContainerPort int

var organizationRefreshToken string

var privacyAPIHost string
var privacyAPIScheme string
var privacyContainerPort int

var sorID string
var sorURL string
var sorOrganizationCode string

var vaultAPIHost string
var vaultAPIScheme string
var vaultContainerPort int
var vaultRefreshToken string
var vaultSealUnsealKey string

var azureServiceBusConnectionString string

var sapAPIHost string
var sapAPIScheme string
var sapAPIUsername string
var sapAPIPassword string
var sapAPIPath string

var serviceNowAPIHost string
var serviceNowAPIScheme string
var serviceNowAPIUsername string
var serviceNowAPIPassword string
var serviceNowAPIPath string

var salesforceAPIHost string
var salesforceAPIScheme string
var salesforceAPIPath string

var withLocalVault bool
var withLocalIdent bool
var withLocalNChain bool
var withLocalPrivacy bool

var withoutPrivateSubnet bool
var withoutRequireOrganizationKeys bool
var withoutRequireSubjectAccount bool
var withoutRequireWorkgroup bool

var withLocalBaselineBuild bool
var baselineDirPath string

var startBaselineStackCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the baseline stack",
	Long:  `Start a local BPI stack instance and connect to internal systems of record`,
	Run:   startStack,
}

func startStack(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepStart)
}

func runStackStart(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	go common.PurgeContainers(docker, name, prune)

	authorizeContext()
	normalizeHostnames()
	sorPrompt()
	tunnelAPIPrompt()
	tunnelMessagingPrompt()

	wg := &sync.WaitGroup{}

	images := make([]string, 0)
	images = append(
		images,
		baselineContainerImage,
		elasticContainerImage,
		natsContainerImage,
		postgresContainerImage,
		redisContainerImage,
	)

	if withLocalIdent {
		identVersion := "latest"
		if common.IsReleaseContext() {
			version, err := common.Manifest.GetImageVersion(identContainerImage)
			if err != nil {
				log.Printf("failed to resolve version for pinned container image: %s; %s", identContainerImage, err.Error())
				os.Exit(1)
			}
			identVersion = *version
		}
		identImage := fmt.Sprintf("%s:%s", identContainerImage, identVersion)
		images = append(images, identImage)
	} else if common.IsReleaseContext() {
		// TODO: enforce target API major version level using status endpoint at configured identBaseURL/status
	}

	if withLocalNChain {
		nchainVersion := "latest"
		if common.IsReleaseContext() {
			version, err := common.Manifest.GetImageVersion(nchainContainerImage)
			if err != nil {
				log.Printf("failed to resolve version for pinned container image: %s; %s", nchainContainerImage, err.Error())
				os.Exit(1)
			}
			nchainVersion = *version
		}
		nchainImage := fmt.Sprintf("%s:%s", nchainContainerImage, nchainVersion)
		images = append(images, nchainImage)
	} else if common.IsReleaseContext() {
		// TODO: enforce target API major version level using status endpoint at configured nchainBaseURL/status
	}

	if withLocalPrivacy {
		privacyVersion := "latest"
		if common.IsReleaseContext() {
			version, err := common.Manifest.GetImageVersion(privacyContainerImage)
			if err != nil {
				log.Printf("failed to resolve version for pinned container image: %s; %s", privacyContainerImage, err.Error())
				os.Exit(1)
			}
			privacyVersion = *version
		}
		privacyImage := fmt.Sprintf("%s:%s", privacyContainerImage, privacyVersion)
		images = append(images, privacyImage)
	} else if common.IsReleaseContext() {
		// TODO: enforce target API major version level using status endpoint at configured privacyBaseURL/status
	}

	if withLocalVault {
		vaultVersion := "latest"
		if common.IsReleaseContext() {
			version, err := common.Manifest.GetImageVersion(vaultContainerImage)
			if err != nil {
				log.Printf("failed to resolve version for pinned container image: %s; %s", vaultContainerImage, err.Error())
				os.Exit(1)
			}
			vaultVersion = *version
		}
		vaultImage := fmt.Sprintf("%s:%s", vaultContainerImage, vaultVersion)
		images = append(images, vaultImage)
	} else if common.IsReleaseContext() {
		// TODO: enforce target API major version level using status endpoint at configured vaultBaseURL/status
	}

	if withLocalIdent || withLocalNChain || withLocalPrivacy || withLocalVault {
		images = append(images, postgresContainerImage)
	}

	for _, image := range images {
		img := image
		wg.Add(1)
		go func() {
			err := pullImage(docker, img)
			if err != nil {
				log.Printf("failed to pull local BPI container image: %s; %s", img, err.Error())
				os.Exit(1)
			}
			wg.Done()
		}()
	}

	configureNetwork(docker)
	common.RequireOrganizationEndpoints(
		func() {
			applyFlags()

			wg.Wait()

			// run local deps
			wg.Add(1)
			go runElasticsearch(docker, wg)

			// FIXME-- DRY this up...
			elasticReachable := false
			for !elasticReachable {
				host := fmt.Sprintf("localhost:%v", elasticPort)
				conn, err := net.DialTimeout("tcp", host, defaultElasticReachabilityTimeout)
				if err == nil {
					conn.Close()
					elasticReachable = true
				}
			}

			wg.Add(1)
			go runNATS(docker, wg)

			// FIXME-- DRY this up...
			natsReachable := false
			for !natsReachable {
				host := fmt.Sprintf("localhost:%v", natsPort)
				conn, err := net.DialTimeout("tcp", host, defaultNatsReachabilityTimeout)
				if err == nil {
					conn.Close()
					natsReachable = true
				}
			}

			wg.Add(1)
			go runPostgres(docker, wg)

			// FIXME-- DRY this up...
			postgresReachable := false
			for !postgresReachable {
				host := fmt.Sprintf("localhost:%v", postgresPort)
				conn, err := net.DialTimeout("tcp", host, defaultPostgresReachabilityTimeout)
				if err == nil {
					conn.Close()
					postgresReachable = true
				}
			}

			wg.Add(1)
			go runRedis(docker, wg)

			// FIXME-- DRY this up...
			redisReachable := false
			for !redisReachable {
				host := fmt.Sprintf("localhost:%v", redisPort)
				conn, err := net.DialTimeout("tcp", host, defaultRedisReachabilityTimeout)
				if err == nil {
					conn.Close()
					redisReachable = true
				}
			}

			// run optional local containers
			if withLocalIdent {
				wg.Add(1)
				go runIdentAPI(docker, wg)

				wg.Add(1)
				go runIdentConsumer(docker, wg)
			}

			if withLocalNChain {
				wg.Add(1)
				go runNChainAPI(docker, wg)

				wg.Add(1)
				go runNChainConsumer(docker, wg)

				wg.Add(1)
				go runStatsdaemon(docker, wg)

				wg.Add(1)
				go runReachabilitydaemon(docker, wg)
			}

			if withLocalPrivacy {
				wg.Add(1)
				go runPrivacyAPI(docker, wg)

				wg.Add(1)
				go runPrivacyConsumer(docker, wg)
			}

			if withLocalVault {
				wg.Add(1)
				go runVaultAPI(docker, wg)
			}

			if withLocalBaselineBuild {
				useLocalBaselineBuild(docker)
			}

			// run BPI
			wg.Add(1)
			go runBaselineAPI(docker, wg)

			wg.Add(1)
			go runBaselineConsumer(docker, wg)

			wg.Wait()
			log.Printf("%s local BPI instance started", name)

			if !withoutRequireOrganizationKeys {
				requireOrganizationKeys()
			}

			if !withoutRequireSubjectAccount {
				requireBPISubjectAccount()
			}
		},
		func(reason *string) {
			if reason != nil {
				log.Printf(*reason)

				if !prune {
					common.StopContainers(docker, name)
				} else {
					common.PurgeContainers(docker, name, true)
					common.PurgeNetwork(docker, name)
				}
			}
		},
		port,
		natsPort,
		natsWebsocketPort,
	)
}

func sorConfigFactory() map[string]interface{} {
	sor := map[string]interface{}{}
	// TODO-- write the SOR configuration...

	// fmt.Sprintf("PROVIDE_SOR_IDENTIFIER=%s", sorID),
	// fmt.Sprintf("PROVIDE_SOR_ORGANIZATION_CODE=%s", sorOrganizationCode),
	// fmt.Sprintf("PROVIDE_SOR_URL=%s", sorURL),

	// if azureServiceBusConnectionString != "" {
	// 	for _, envvar := range []string{
	// 		fmt.Sprintf("AZURE_SERVICE_BUS_CONNECTION_STRING=%s", azureServiceBusConnectionString),
	// 	} {
	// 		env = append(env, envvar)
	// 	}
	// }

	// if sapAPIHost != "" && sapAPIUsername != "" && sapAPIPassword != "" {
	// 	for _, envvar := range []string{
	// 		fmt.Sprintf("SAP_API_HOST=%s", sapAPIHost),
	// 		fmt.Sprintf("SAP_API_SCHEME=%s", sapAPIScheme),
	// 		fmt.Sprintf("SAP_API_PATH=%s", sapAPIPath),
	// 		fmt.Sprintf("SAP_API_USERNAME=%s", sapAPIUsername),
	// 		fmt.Sprintf("SAP_API_PASSWORD=%s", sapAPIPassword),
	// 	} {
	// 		env = append(env, envvar)
	// 	}
	// } else if sorID == "sap" && sorURL != "" {
	// 	_url, err := url.Parse(sorURL)
	// 	if err != nil {
	// 		log.Printf("WARNING: system of record url invalid; %s", err.Error())
	// 	}
	// 	for _, envvar := range []string{
	// 		fmt.Sprintf("SAP_API_HOST=%s", _url.Host),
	// 		fmt.Sprintf("SAP_API_SCHEME=%s", _url.Scheme),
	// 	} {
	// 		env = append(env, envvar)
	// 	}

	// 	if _url.Path != "" {
	// 		env = append(env, fmt.Sprintf("SAP_API_PATH=%s", strings.TrimLeft(_url.Path, "/")))
	// 	}
	// }

	// if serviceNowAPIHost != "" && serviceNowAPIUsername != "" && serviceNowAPIPassword != "" {
	// 	for _, envvar := range []string{
	// 		fmt.Sprintf("SERVICENOW_API_HOST=%s", serviceNowAPIHost),
	// 		fmt.Sprintf("SERVICENOW_API_SCHEME=%s", serviceNowAPIScheme),
	// 		fmt.Sprintf("SERVICENOW_API_PATH=%s", serviceNowAPIPath),
	// 		fmt.Sprintf("SERVICENOW_API_USERNAME=%s", serviceNowAPIUsername),
	// 		fmt.Sprintf("SERVICENOW_API_PASSWORD=%s", serviceNowAPIPassword),
	// 	} {
	// 		env = append(env, envvar)
	// 	}
	// } else if sorID == "servicenow" || sorID == "snow" && sorURL != "" {
	// 	_url, err := url.Parse(sorURL)
	// 	if err != nil {
	// 		log.Printf("WARNING: system of record url invalid; %s", err.Error())
	// 	}
	// 	for _, envvar := range []string{
	// 		fmt.Sprintf("SERVICENOW_API_HOST=%s", _url.Host),
	// 		fmt.Sprintf("SERVICENOW_API_SCHEME=%s", _url.Scheme),
	// 	} {
	// 		env = append(env, envvar)
	// 	}

	// 	if _url.Path != "" {
	// 		env = append(env, fmt.Sprintf("SERVICENOW_API_PATH=%s", strings.TrimLeft(_url.Path, "/")))
	// 	}
	// }

	return sor
}

func requireBPISubjectAccount() error {
	log.Printf("waiting for BPI to become available...")
	for baseline.Status() != nil {
		time.Sleep(time.Second * 1)
	}
	log.Printf("BPI is available")

	token, err := ident.CreateToken(organizationRefreshToken, map[string]interface{}{
		"grant_type":      "refresh_token",
		"organization_id": common.OrganizationID,
	})
	if err != nil {
		log.Printf("failed to authorize access token on behalf of organization %s; %s", common.OrganizationID, err.Error())
		os.Exit(1)
	}

	var sacct *baseline.SubjectAccount
	subjectAccountID := baseline.SubjectAccountIDFactory(common.OrganizationID, common.WorkgroupID)

	sacct, err = baseline.GetSubjectAccountDetails(*token.AccessToken, common.OrganizationID, subjectAccountID, map[string]interface{}{})
	if err == nil && sacct != nil && sacct.ID != nil {
		log.Printf("BPI subject account resolved: %s", *sacct.ID)
		// TODO-- update if needed...
		return nil
	}

	sacct, err = baseline.CreateSubjectAccount(*token.AccessToken, common.OrganizationID, map[string]interface{}{
		"metadata": &baseline.SubjectAccountMetadata{
			// Counterparties []*Participant `sql:"-" json:"counterparties,omitempty"`
			NetworkID:           &nchainBaselineNetworkID,
			OrganizationAddress: &baselineOrganizationAddress,
			// OrganizationDomain *string `json:"organization_domain,omitempty"`
			OrganizationID:                &common.OrganizationID,
			OrganizationMessagingEndpoint: &common.MessagingEndpoint,
			OrganizationProxyEndpoint:     &common.BPIEndpoint,
			OrganizationRefreshToken:      &organizationRefreshToken,
			// OrganizationWebsocketEndpoint *string `json:"organization_websocket_endpoint,omitempty"`
			RegistryContractAddress: &baselineRegistryContractAddress,
			SOR:                     sorConfigFactory(),
			WorkgroupID:             &common.WorkgroupID,
		},
		"subject_id": common.OrganizationID,
	})

	if err != nil {
		log.Printf("WARNING: BPI subject account not created; %s", err.Error())
		return err
	}

	log.Printf("BPI subject account created: %s", *sacct.ID)
	return nil
}

func requireOrganizationKeys() {
	var err error

	_, err = common.RequireOrganizationKeypair("babyJubJub")
	if err != nil {
		log.Printf("WARNING: failed to require organization keypair; %s", err.Error())
	}

	_, err = common.RequireOrganizationKeypair("secp256k1")
	if err != nil {
		log.Printf("WARNING: failed to require organization keypair; %s", err.Error())
	}

	_, err = common.RequireOrganizationKeypair("BIP39")
	if err != nil {
		log.Printf("WARNING: failed to require organization keypair; %s", err.Error())
	}

	_, err = common.RequireOrganizationKeypair("RSA-4096")
	if err != nil {
		log.Printf("WARNING: failed to require organization keypair; %s", err.Error())
	}
}

func configureNetwork(docker *client.Client) {
	if withoutPrivateSubnet {
		log.Printf("skipping initialization of a private subnet")
		return
	}

	opts := types.NetworkCreate{
		// CheckDuplicate bool
		// Driver		  string
		// Scope          string
		// EnableIPv6     bool
		IPAM: &network.IPAM{},
		// Internal       bool
		// Attachable     bool
		// Ingress        bool
		// ConfigOnly     bool
		// ConfigFrom     *network.ConfigReference
		// Options        map[string]string
		// Labels         map[string]string
	}

	if !strings.EqualFold(runtime.GOOS, "windows") {
		opts.Driver = dockerNetworkDriverBridge
	}

	network, err := docker.NetworkCreate(
		context.Background(),
		name,
		opts,
	)

	if err != nil {
		log.Printf("failed to setup docker network; %s", err.Error())
		os.Exit(1)
	}

	dockerNetworkID = network.ID
	log.Printf("configured network for local BPI instance: %s", name)
}

func normalizeHostnames() {
	apiHostname = strings.Replace(apiHostname, defaultBPIStackName, name, -1)
	consumerHostname = strings.Replace(consumerHostname, defaultBPIStackName, name, -1)
	elasticHostname = strings.Replace(elasticHostname, defaultBPIStackName, name, -1)
	identHostname = strings.Replace(identHostname, defaultBPIStackName, name, -1)
	identConsumerHostname = strings.Replace(identConsumerHostname, defaultBPIStackName, name, -1)
	natsHostname = strings.Replace(natsHostname, defaultBPIStackName, name, -1)
	nchainHostname = strings.Replace(nchainHostname, defaultBPIStackName, name, -1)
	nchainConsumerHostname = strings.Replace(nchainConsumerHostname, defaultBPIStackName, name, -1)
	postgresHostname = strings.Replace(postgresHostname, defaultBPIStackName, name, -1)
	privacyHostname = strings.Replace(privacyHostname, defaultBPIStackName, name, -1)
	privacyConsumerHostname = strings.Replace(privacyConsumerHostname, defaultBPIStackName, name, -1)
	redisHostname = strings.Replace(redisHostname, defaultBPIStackName, name, -1)
	vaultHostname = strings.Replace(vaultHostname, defaultBPIStackName, name, -1)
}

func authorizeContext() {
	if !withoutRequireWorkgroup {
		// log.Printf("authorizing workgroup context")
		authorizeWorkgroupContext()
	}

	// log.Printf("authorizing organization context")
	common.AuthorizeOrganizationContext(false)

	if organizationRefreshToken == "" {
		refreshTokenKey := common.BuildConfigKeyWithID(common.RefreshTokenConfigKey, common.OrganizationID)
		if viper.IsSet(refreshTokenKey) {
			// log.Printf("using cached API refresh token for organization: %s\n", common.OrganizationID)
			organizationRefreshToken = viper.GetString(refreshTokenKey)
			if vaultRefreshToken == "" {
				vaultRefreshToken = organizationRefreshToken
			}
		} else {
			organizationAuthPrompt()
			if common.OrganizationRefreshToken != "" {
				organizationRefreshToken = common.OrganizationRefreshToken
				if vaultRefreshToken == "" {
					vaultRefreshToken = organizationRefreshToken
				}
			} else {
				log.Printf("failed to resolve refresh token for organization: %s\n", common.OrganizationID)
				os.Exit(1)
			}
		}
	} else if organizationRefreshToken != "" && vaultRefreshToken == "" {
		vaultRefreshToken = organizationRefreshToken
	}
}

func authorizeWorkgroupContext() {
	if common.WorkgroupID == "" {
		err := common.RequireWorkgroup()
		if err != nil {
			log.Printf("failed to require workgroup; %s", err.Error())
			os.Exit(1)
		}
	}

	var contracts []*nchain.Contract

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		log.Printf("failed to resolve workgroup: %s; %s", common.WorkgroupID, err.Error())
		os.Exit(1)
	}

	workgroup, err := ident.GetApplicationDetails(*token.AccessToken, common.WorkgroupID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to resolve workgroup: %s; %s", common.WorkgroupID, err.Error())
		os.Exit(1)
	}

	contracts, err = nchain.ListContracts(*token.AccessToken, map[string]interface{}{
		"type": "organization-registry",
	})
	if err != nil {
		log.Printf("failed to resolve global organization registry contract; %s", err.Error())
		os.Exit(1)
	} else if len(contracts) == 0 {
		common.AuthorizeOrganizationContext(true)

		token, err := ident.CreateToken(*token.AccessToken, map[string]interface{}{
			"scope":           "offline_access",
			"organization_id": common.OrganizationID,
		})
		if err != nil {
			log.Printf("failed to authorize API access token on behalf of workgroup %s; %s", common.WorkgroupID, err.Error())
			os.Exit(1)
		}

		contracts, err = nchain.ListContracts(*token.AccessToken, map[string]interface{}{
			"type": "organization-registry",
		})
		if err != nil {
			log.Printf("failed to resolve global organization registry contract; %s", err.Error())
			os.Exit(1)
		} else if len(contracts) == 0 {
			log.Printf("failed to resolve global organization registry contract")
			os.Exit(1)
		}
	}

	if nchainBaselineNetworkID == "" {
		if workgroup.NetworkID != uuid.Nil {
			nchainBaselineNetworkID = workgroup.NetworkID.String()
		} else {
			err := common.RequireL1Network()
			if err != nil {
				log.Printf("failed to require network id; %s", err.Error())
				os.Exit(1)
			}
			nchainBaselineNetworkID = common.NetworkID
		}
	}

	orgRegistryContract := contracts[0]
	if orgRegistryContract.Address == nil || *orgRegistryContract.Address == "0x" {
		log.Printf("failed to resolve global organization registry contract; %s", err.Error())
		os.Exit(1)
	}
	baselineRegistryContractAddress = *contracts[0].Address
}

func applyFlags() {
	if (baselineOrganizationAddress == "" || baselineOrganizationAddress == "0x") && common.ResolvedBaselineOrgAddress != "" {
		// FIXME-- this belongs somewhere better
		baselineOrganizationAddress = common.ResolvedBaselineOrgAddress
	}

	// HACK
	if strings.HasSuffix(apiHostname, "-api") {
		apiHostname = fmt.Sprintf("%s-api", name)
	}

	// HACK
	if strings.HasSuffix(consumerHostname, "-consumer") {
		consumerHostname = fmt.Sprintf("%s-consumer", name)
	}

	// HACK
	if strings.HasSuffix(identHostname, "-ident-api") {
		identHostname = fmt.Sprintf("%s-ident-api", name)
	}

	// HACK
	if strings.HasSuffix(identHostname, "-ident-consumer") {
		identConsumerHostname = fmt.Sprintf("%s-ident-consumer", name)
	}

	// HACK
	if strings.HasSuffix(nchainHostname, "-nchain-api") {
		nchainHostname = fmt.Sprintf("%s-nchain-api", name)
	}

	// HACK
	if strings.HasSuffix(nchainConsumerHostname, "-nchain-consumer") {
		nchainConsumerHostname = fmt.Sprintf("%s-nchain-consumer", name)
	}

	// HACK
	if strings.HasSuffix(nchainConsumerHostname, "-reachabilitydaemon") {
		nchainReachabilitydaemonHostname = fmt.Sprintf("%s-reachabilitydaemon", name)
	}

	// HACK
	if strings.HasSuffix(nchainStatsdaemonHostname, "-statsdaemon") {
		nchainStatsdaemonHostname = fmt.Sprintf("%s-statsdaemon", name)
	}

	// HACK
	if strings.HasSuffix(privacyHostname, "-privacy-api") {
		privacyHostname = fmt.Sprintf("%s-privacy-api", name)
	}

	// HACK
	if strings.HasSuffix(privacyConsumerHostname, "-privacy-consumer") {
		privacyConsumerHostname = fmt.Sprintf("%s-privacy-consumer", name)
	}

	// HACK
	if strings.HasSuffix(vaultHostname, "-vault-api") {
		vaultHostname = fmt.Sprintf("%s-vault-api", name)
	}

	// HACK
	if strings.HasSuffix(natsHostname, "-nats") {
		natsHostname = fmt.Sprintf("%s-nats", name)
	}

	// HACK
	if strings.HasSuffix(redisHostname, "-redis") {
		redisHostname = fmt.Sprintf("%s-redis", name)
		redisHosts = fmt.Sprintf("%s:%d", redisHostname, redisContainerPort)
	}

	// HACK
	if natsServerName == "" {
		natsServerName = defaultNatsServerName
	}

	// HACK
	if jwtSignerPublicKey == "" {
		keys, err := vault.ListKeys(common.OrganizationAccessToken, common.VaultID, map[string]interface{}{
			"spec": "RSA-4096",
		})
		if err != nil {
			log.Printf("WARNING: failed to resolve RSA-4096 key for organization; %s", err.Error())
			return
		}
		if len(keys) == 0 {
			log.Printf("WARNING: failed to resolve RSA-4096 key for organization")
			return
		}

		jwtSignerPublicKey = *keys[0].PublicKey
	}
}

func containerEnvironmentFactory(listenPort *int) []string {
	env := make([]string, 0)
	for _, envvar := range []string{
		fmt.Sprintf("DATABASE_HOST=%s", postgresHostname),
		fmt.Sprintf("DATABASE_PORT=%d", postgresContainerPort),
		fmt.Sprintf("DATABASE_USER=%s", postgresUser),
		fmt.Sprintf("DATABASE_PASSWORD=%s", postgresPassword),
		fmt.Sprintf("DATABASE_NAME=%s", postgresDatabase),
		fmt.Sprintf("DATABASE_SUPERUSER=%s", "prvd"),
		fmt.Sprintf("DATABASE_SUPERUSER_PASSWORD=%s", "prvdp455"),
		fmt.Sprintf("DATABASE_LOGGING=%s", databaseLogging),
		fmt.Sprintf("ELASTICSEARCH_ACCEPT_SELF_SIGNED_CERTIFICATE=%t", elasticAcceptSelfSignedCertificate),
		fmt.Sprintf("ELASTICSEARCH_API_SCHEME=%s", elasticAPIScheme),
		fmt.Sprintf("ELASTICSEARCH_HOSTS=%s:%d", elasticHostname, elasticContainerPort),
		fmt.Sprintf("ELASTICSEARCH_USERNAME=%s", elasticUsername),
		fmt.Sprintf("ELASTICSEARCH_PASSWORD=%s", elasticPassword),
		fmt.Sprintf("IDENT_API_HOST=%s", identAPIHost),
		fmt.Sprintf("IDENT_API_SCHEME=%s", identAPIScheme),
		fmt.Sprintf("JWT_SIGNER_PUBLIC_KEY=%s", strings.ReplaceAll(jwtSignerPublicKey, "\\n", "\n")),
		fmt.Sprintf("LOG_LEVEL=%s", logLevel),
		fmt.Sprintf("NATS_CLIENT_PREFIX=%s", name),
		fmt.Sprintf("NATS_JETSTREAM_URL=%s", fmt.Sprintf("nats://%s:%d", natsHostname, natsContainerPort)),
		fmt.Sprintf("NATS_TOKEN=%s", natsAuthToken),
		fmt.Sprintf("NATS_URL=%s", fmt.Sprintf("nats://%s:%d", natsHostname, natsContainerPort)),
		fmt.Sprintf("NCHAIN_API_HOST=%s", nchainAPIHost),
		fmt.Sprintf("NCHAIN_API_SCHEME=%s", nchainAPIScheme),
		fmt.Sprintf("NCHAIN_BASELINE_NETWORK_ID=%s", nchainBaselineNetworkID),
		fmt.Sprintf("PRIVACY_API_HOST=%s", privacyAPIHost),
		fmt.Sprintf("PRIVACY_API_SCHEME=%s", privacyAPIScheme),
		fmt.Sprintf("PRIVACY_API_SCHEME=%s", privacyAPIScheme),
		fmt.Sprintf("REDIS_HOSTS=%s", redisHosts),
		fmt.Sprintf("SYSLOG_ENDPOINT=%s", syslogEndpoint),
		fmt.Sprintf("VAULT_API_HOST=%s", vaultAPIHost),
		fmt.Sprintf("VAULT_API_SCHEME=%s", vaultAPIScheme),
		fmt.Sprintf("VAULT_REFRESH_TOKEN=%s", vaultRefreshToken),
		fmt.Sprintf("VAULT_SEAL_UNSEAL_KEY=%s", vaultSealUnsealKey),

		"SERVICENOW_LIST_SCHEMAS_API_PATH=api/x_prot9_provide/setup_configuration/ListSchemas",
		"SERVICENOW_SCHEMA_DETAILS_API_PATH=api/x_prot9_provide/setup_configuration/GetSchemaDetails",
		"SERVICENOW_HEALTHCHECK_API_PATH=api/now/branding",
	} {
		env = append(env, envvar)
	}

	if listenPort != nil {
		env = append(env, fmt.Sprintf("PORT=%d", *listenPort))
	}

	return env
}

func runBaselineAPI(docker *client.Client, wg *sync.WaitGroup) {
	image := baselineContainerImage
	if withLocalBaselineBuild {
		image = localBaselineContainerImage
	}

	err := runContainer(
		docker,
		fmt.Sprintf("%s-api", strings.ReplaceAll(name, " ", "")),
		apiHostname,
		image,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", apiHostname, containerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      port,
			containerPort: containerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local BPI container; %s", err.Error())
		os.Exit(1)
	}

	os.Setenv("BASELINE_API_HOST", fmt.Sprintf("localhost:%d", port))
	os.Setenv("BASELINE_API_SCHEME", "http")

	if wg != nil {
		wg.Done()
	}
}

func runBaselineConsumer(docker *client.Client, wg *sync.WaitGroup) {
	image := baselineContainerImage
	if withLocalBaselineBuild {
		image = localBaselineContainerImage
	}

	err := runContainer(
		docker,
		fmt.Sprintf("%s-consumer", strings.ReplaceAll(name, " ", "")),
		consumerHostname,
		image,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", apiHostname, port)},
		nil,
		map[string]string{},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create local BPI consumer container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runIdentAPI(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-ident-api", strings.ReplaceAll(name, " ", "")),
		identHostname,
		identContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", identHostname, identContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      identPort,
			containerPort: identContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local ident API container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runIdentConsumer(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-ident-consumer", strings.ReplaceAll(name, " ", "")),
		identConsumerHostname,
		identContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", identHostname, identContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create local ident consumer container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runNChainAPI(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-nchain-api", strings.ReplaceAll(name, " ", "")),
		nchainHostname,
		nchainContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", nchainHostname, nchainContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      nchainPort,
			containerPort: nchainContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local nchain API container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runNChainConsumer(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-nchain-consumer", strings.ReplaceAll(name, " ", "")),
		nchainConsumerHostname,
		nchainContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", nchainHostname, nchainContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create local nchain consumer container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runStatsdaemon(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-statsdaemon", strings.ReplaceAll(name, " ", "")),
		nchainStatsdaemonHostname,
		nchainContainerImage,
		&[]string{"./ops/run_statsdaemon.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", nchainHostname, nchainContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create local statsdaemon container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runReachabilitydaemon(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-reachabilitydaemon", strings.ReplaceAll(name, " ", "")),
		nchainReachabilitydaemonHostname,
		nchainContainerImage,
		&[]string{"./ops/run_reachabilitydaemon.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", nchainHostname, nchainContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create local reachabilitydaemon container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runPrivacyAPI(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-privacy-api", strings.ReplaceAll(name, " ", "")),
		privacyHostname,
		privacyContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", privacyHostname, privacyContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      privacyPort,
			containerPort: privacyContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local privacy API container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runPrivacyConsumer(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-privacy-consumer", strings.ReplaceAll(name, " ", "")),
		privacyConsumerHostname,
		privacyContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", privacyHostname, privacyContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create local privacy consumer container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runVaultAPI(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-vault-api", strings.ReplaceAll(name, " ", "")),
		vaultHostname,
		vaultContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", vaultHostname, vaultContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      vaultPort,
			containerPort: vaultContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local vault API container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func writeNATSConfig() *string {
	cfg := []byte("max_payload: 100Mb\nmax_pending: 104857600\n")
	if !natsWebsocketTLS {
		cfg = []byte(fmt.Sprintf("max_payload: 100Mb\nmax_pending: 104857600\nwebsocket {\n    listen: \"0.0.0.0:%d\"\n    no_tls: true\n}\n", natsWebsocketContainerPort))
	}
	path := strings.Split(os.TempDir(), string(os.PathSeparator))
	path = append(path, "nats-server.conf")
	sep := []string{string(os.PathSeparator)}
	path = append(sep, path...)
	tmp := filepath.Join(path...)
	err := ioutil.WriteFile(filepath.FromSlash(strings.ReplaceAll(tmp, string(os.PathSeparator), "/")), cfg, 0644)
	if err != nil {
		log.Printf("failed to write local nats-server.conf; %s", err.Error())
		os.Exit(1)
	}

	if tmp == "" {
		return nil
	}
	return &tmp
}

func runElasticsearch(docker *client.Client, wg *sync.WaitGroup) {
	// cfgPath := writeElasticsearchConfig()
	mountPoints := map[string]string{}

	// if cfgPath != nil {
	// 	mountPoints[*cfgPath] = "/etc/elasticsearch.conf"
	// }

	env := []string{
		fmt.Sprintf("ELASTIC_PASSWORD=%s", elasticPassword),
		"discovery.type=single-node",
		"bootstrap.memory_lock=true",
		// "-p", fmt.Sprintf("%d", elasticContainerPort),
		// "--ulimit", "nofile=65535:65535",
		// "http.host=0.0.0.0",
		// "transport.host=0.0.0.0",
		// "xpack.security.enabled=false",
		// "xpack.security.http.ssl.enabled=false",
	}

	if elasticMemory != "" {
		_elasticMemory, err := strconv.ParseUint(strings.Replace(elasticMemory, "GB", "", -1), 10, 64)
		if err != nil {
			log.Printf("failed to parse provided elasticsearch memory allocation; %s", err.Error())
			os.Exit(1)
		}
		env = append(env, fmt.Sprintf("ES_JAVA_OPTS=-Xms%dm -Xmx%dm", _elasticMemory*1024, _elasticMemory*1024))
	}

	err := runContainer(
		docker,
		fmt.Sprintf("%s-elasticsearch", strings.ReplaceAll(name, " ", "")),
		elasticHostname,
		elasticContainerImage,
		nil,
		nil,
		&[]string{"CMD", "nc", "-zv", "localhost", fmt.Sprintf("%d", elasticPort)},
		&env,
		mountPoints,
		[]portMapping{
			{
				hostPort:      elasticPort,
				containerPort: elasticContainerPort,
			},
		}...,
	)

	if err != nil {
		log.Printf("failed to create local BPI elasticsearch container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runNATS(docker *client.Client, wg *sync.WaitGroup) {
	cfgPath := writeNATSConfig()
	mountPoints := map[string]string{}

	if cfgPath != nil {
		mountPoints[*cfgPath] = "/etc/nats-server.conf"
	}

	err := runContainer(
		docker,
		fmt.Sprintf("%s-nats", strings.ReplaceAll(name, " ", "")),
		natsHostname,
		natsContainerImage,
		nil,
		&[]string{
			"--js",
			"--server_name", natsServerName,
			"--auth", natsAuthToken,
			"--config", "/etc/nats-server.conf",
			"--port", fmt.Sprintf("%d", natsContainerPort),
			"-DVV",
		},
		&[]string{"CMD", "nc", "-zv", "localhost", fmt.Sprintf("%d", natsContainerPort)},
		nil,
		mountPoints,
		[]portMapping{
			{
				hostPort:      natsPort,
				containerPort: natsContainerPort,
			},
			{
				hostPort:      natsWebsocketPort,
				containerPort: natsWebsocketContainerPort,
			},
		}...,
	)

	if err != nil {
		log.Printf("failed to create local BPI NATS container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runPostgres(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-postgres", strings.ReplaceAll(name, " ", "")),
		postgresHostname,
		postgresContainerImage,
		nil,
		nil,
		&[]string{"CMD", "pg_isready", "-U", "prvd", "-d", "prvd"},
		&[]string{
			// FIXME -- allow user to set these....
			"POSTGRES_DB=prvd",
			"POSTGRES_USER=prvd",
			"POSTGRES_PASSWORD=prvdp455",
		},
		map[string]string{},
		[]portMapping{{
			hostPort:      postgresPort,
			containerPort: postgresContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local postgres container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runRedis(docker *client.Client, wg *sync.WaitGroup) {
	err := runContainer(
		docker,
		fmt.Sprintf("%s-redis", strings.ReplaceAll(name, " ", "")),
		redisHostname,
		redisContainerImage,
		nil,
		nil,
		&[]string{"CMD", "redis-cli", "ping"},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      redisPort,
			containerPort: redisContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local BPI redis container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func pullImage(docker *client.Client, image string) error {
	log.Printf("pulling local BPI container image: %s", image)
	reader, err := docker.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("WARNING: %s", err.Error())
	}

	// log.Printf("%s", string(buf))
	io.Copy(os.Stdout, reader)

	return nil
}

func runContainer(
	docker *client.Client,
	name, hostname, image string,
	entrypoint, cmd, healthcheck, env *[]string,
	mounts map[string]string,
	ports ...portMapping,
) error {
	log.Printf("running local BPI container image: %s", image)

	isReachable := func(host string, port int) bool {
		addr := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", addr, defaultContainerReachabilityTimeout)
		if err == nil {
			defer conn.Close()
		}
		return err == nil
	}

	portBinding := nat.PortMap{}
	for _, mapping := range ports {
		if isReachable("0.0.0.0", mapping.hostPort) {
			log.Printf("failed to run local BPI container image: %s; bind for 0.0.0.0:%d failed; port is already allocated", image, mapping.hostPort)
			os.Exit(1)
		}

		port, _ := nat.NewPort("tcp", strconv.Itoa(mapping.containerPort))
		portBinding[port] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: strconv.Itoa(mapping.hostPort),
		}}
	}

	var listenPort *int
	if len(ports) == 1 {
		listenPort = &ports[0].containerPort
	}

	var environment []string
	if env != nil {
		environment = *env
	} else {
		environment = containerEnvironmentFactory(listenPort)
	}

	containerConfig := &container.Config{
		Env:      environment,
		Hostname: hostname,
		Image:    image,
	}

	if cmd != nil {
		containerConfig.Cmd = *cmd
	}

	if entrypoint != nil {
		containerConfig.Entrypoint = *entrypoint
	}

	if healthcheck != nil {
		containerConfig.Healthcheck = &container.HealthConfig{
			Interval:    time.Minute * 1,
			Retries:     2,
			StartPeriod: time.Second * 10,
			Test:        *healthcheck,
			Timeout:     time.Second * 1,
		}
	}

	mountedVolumes := make([]mount.Mount, 0)
	for source := range mounts { // mounts are mapped source => target...
		mountedVolumes = append(mountedVolumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: mounts[source],
		})
	}

	var containerID string
	for _, container := range common.ListContainers(docker, "") {
		if strings.ReplaceAll(container.Names[0], "/", "") == name {
			containerID = container.ID
		}
	}

	if containerID == "" {
		hostConfig := &container.HostConfig{
			AutoRemove:   autoRemove,
			ExtraHosts:   []string{"host.docker.internal:host-gateway"},
			Mounts:       mountedVolumes,
			NetworkMode:  dockerNetworkDriverDefault,
			PortBindings: portBinding,
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		}

		if !strings.EqualFold(runtime.GOOS, "windows") {
			hostConfig.NetworkMode = dockerNetworkDriverBridge
		}

		container, err := docker.ContainerCreate(
			context.Background(),
			containerConfig,
			hostConfig,
			&network.NetworkingConfig{},
			nil,
			strings.ReplaceAll(name, " ", ""),
		)

		if err != nil {
			return err
		}

		containerID = container.ID
	}

	err := docker.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	err = docker.NetworkConnect(
		context.Background(),
		dockerNetworkID,
		containerID,
		&network.EndpointSettings{},
	)
	if err != nil {
		return err
	}

	return nil
}

func useLocalBaselineBuild(docker *client.Client) {
	filePath, _ := homedir.Expand(baselineDirPath)
	fileCtx, err := archive.TarWithOptions(filePath, &archive.TarOptions{})
	if err != nil {
		panic(err)
	}

	imageBuildResponse, err := docker.ImageBuild(
		context.Background(),
		fileCtx,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{"baseline_dev"},
		})
	if err != nil {
		log.Fatalf("failed to build image; %s", err.Error())
	}
	defer imageBuildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body) // TODO-- pretty print this
	if err != nil {
		log.Fatalf("failed to read image build response; %s", err.Error())
	}
}

func organizationAuthPrompt() {
	prompt := promptui.Prompt{
		IsConfirm: true,
		Label:     fmt.Sprintf("Authorize access/refresh token for %s", *common.Organization.Name),
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	if strings.ToLower(result) == "y" {
		common.AuthorizeOrganizationContext(true)
	}
}

func tunnelAPIPrompt() {
	if common.WithoutTunnels || common.ExposeBPITunnel || common.BPIEndpoint != "" {
		return
	}

	prompt := promptui.Prompt{
		IsConfirm: true,
		Label:     "Expose tunnel for the local API",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	if strings.ToLower(result) == "y" {
		common.ExposeBPITunnel = true
	}
}

func tunnelMessagingPrompt() {
	if common.WithoutTunnels || common.ExposeMessagingTunnel || common.MessagingEndpoint != "" {
		return
	}

	prompt := promptui.Prompt{
		IsConfirm: true,
		Label:     "Expose tunnel for the local messaging endpoint",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	if strings.ToLower(result) == "y" {
		common.ExposeMessagingTunnel = true
	}
}

func sorPrompt() {
	if sorID != "" {
		return
	}

	items := map[string]string{
		"Dynamics365":           "dynamics365",
		"Ephemeral (In-Memory)": "ephemeral",
		"Excel":                 "excel",
		"Salesforce":            "salesforce",
		"SAP":                   "sap",
		"ServiceNow":            "servicenow",
	}

	opts := make([]string, 0)
	for k := range items {
		opts = append(opts, k)
	}
	sort.Strings(opts)

	prmpt := promptui.Select{
		Label: "What is your primary system of record?",
		Items: opts,
	}

	_, result, _ := prmpt.Run()
	sorID = items[result]

	switch sorID {
	case "ephemeral":
		// no-op
	default:
		sorURLPrompt()
	}
}

func sorURLPrompt() {
	if sorURL != "" {
		return
	}

	prompt := promptui.Prompt{
		Label: "What is the API endpoint for your primary system of record?",
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return
	}

	sorURL = result
}

func init() {
	startBaselineStackCmd.Flags().StringVar(&name, "name", defaultBPIStackName, "name of the baseline stack instance")

	startBaselineStackCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	// runBaselineStackCmd.MarkFlagRequired("organization")

	startBaselineStackCmd.Flags().StringVar(&common.BPIEndpoint, "bpi-endpoint", "", "local BPI endpoint for use by one or more authorized systems of record")
	startBaselineStackCmd.Flags().StringVar(&common.MessagingEndpoint, "messaging-endpoint", "", "public messaging endpoint used for sending and receiving protocol messages")
	startBaselineStackCmd.Flags().BoolVar(&common.Tunnel, "tunnel", false, "when true, a tunnel is established to expose the API and messaging endpoints to the WAN")
	startBaselineStackCmd.Flags().BoolVar(&common.WithoutTunnels, "without-tunnels", false, "when false, no tunnels will be established to expose the API and messaging endpoints to the WAN")
	startBaselineStackCmd.Flags().BoolVar(&common.ExposeBPITunnel, "bpi-tunnel", false, "when true, a tunnel is established to expose the API endpoint to the WAN")
	startBaselineStackCmd.Flags().BoolVar(&common.ExposeMessagingTunnel, "messaging-tunnel", false, "when true, a tunnel is established to expose the messaging endpoint to the WAN")
	startBaselineStackCmd.Flags().BoolVar(&common.ExposeWebsocketMessagingTunnel, "websocket-tunnel", false, "when true, a tunnel is established to expose the websocket messaging endpoint to the WAN")

	startBaselineStackCmd.Flags().StringVar(&sorID, "sor", "", "primary internal system of record identifier being baselined")
	startBaselineStackCmd.Flags().StringVar(&sorURL, "sor-url", "https://", "url of the primary internal system of record being baselined")
	startBaselineStackCmd.Flags().StringVar(&sorOrganizationCode, "sor-organization-code", "", "organization code specific to the system of record")

	startBaselineStackCmd.Flags().StringVar(&apiHostname, "hostname", fmt.Sprintf("%s-api", name), "hostname for the local BPI container")
	startBaselineStackCmd.Flags().IntVar(&port, "port", 8080, "host port on which to expose the local BPI service")
	startBaselineStackCmd.Flags().IntVar(&containerPort, "container-port", defaultBPIContainerPort, "container port on which to expose the BPI service")

	startBaselineStackCmd.Flags().StringVar(&elasticAPIScheme, "elasticsearch-scheme", "https", "protocol scheme of the elasticsearch service")
	startBaselineStackCmd.Flags().StringVar(&elasticHostname, "elasticsearch-hostname", fmt.Sprintf("%s-elasticsearch", name), "hostname for the local BPI elasticsearch container")
	startBaselineStackCmd.Flags().IntVar(&elasticContainerPort, "elasticsearch-container-port", defaultElasticContainerPort, "container port on which to expose the elasticsearch service")
	startBaselineStackCmd.Flags().IntVar(&elasticPort, "elasticsearch-port", 9200, "host port on which to expose the local elasticsearch service")
	startBaselineStackCmd.Flags().StringVar(&elasticUsername, "elasticsearch-username", "elastic", "username of the local elasticsearch service for basic authorization")
	startBaselineStackCmd.Flags().StringVar(&elasticPassword, "elasticsearch-password", "3l4s71c", "password of the local elasticsearch service for basic authorization")
	startBaselineStackCmd.Flags().BoolVar(&elasticAcceptSelfSignedCertificate, "elasticsearch-ssl-insecure", false, "accept self-signed certificate when connecting to the local elasticsearch service")
	startBaselineStackCmd.Flags().StringVar(&elasticMemory, "elasticsearch-memory", "", "amount of RAM allocated to the local elasticsearch service")

	startBaselineStackCmd.Flags().StringVar(&consumerHostname, "consumer-hostname", fmt.Sprintf("%s-consumer", name), "hostname for the local BPI consumer container")
	startBaselineStackCmd.Flags().IntVar(&natsContainerPort, "nats-container-port", defaultNATSContainerPort, "container port on which to expose the NATS service")
	startBaselineStackCmd.Flags().IntVar(&natsWebsocketContainerPort, "nats-ws-container-port", defaultNATSWebsocketContainerPort, "container port on which to expose the NATS websocket service")
	startBaselineStackCmd.Flags().StringVar(&natsHostname, "nats-hostname", fmt.Sprintf("%s-nats", name), "hostname for the local BPI NATS container")
	startBaselineStackCmd.Flags().IntVar(&natsPort, "nats-port", 4222, "host port on which to expose the local NATS service")
	startBaselineStackCmd.Flags().BoolVar(&natsWebsocketTLS, "nats-ws-tls", false, "when true, NATS websocket service uses TLS")
	startBaselineStackCmd.Flags().IntVar(&natsWebsocketPort, "nats-ws-port", 4221, "host port on which to expose the local NATS websocket service")
	startBaselineStackCmd.Flags().StringVar(&natsAuthToken, "nats-auth-token", "testtoken", "authorization token for the local BPI NATS service; will be passed as the -auth argument to NATS")

	startBaselineStackCmd.Flags().IntVar(&postgresContainerPort, "postgres-container-port", defaultPostgresContainerPort, "container port on which to expose the postgres service")
	startBaselineStackCmd.Flags().StringVar(&postgresDatabase, "postgres-database", "baseline", "name for the local postgres database")
	startBaselineStackCmd.Flags().StringVar(&postgresHostname, "postgres-hostname", fmt.Sprintf("%s-postgres", name), "hostname for the local postgres container")
	startBaselineStackCmd.Flags().IntVar(&postgresPort, "postgres-port", 5432, "host port on which to expose the local postgres service")
	startBaselineStackCmd.Flags().StringVar(&postgresUser, "postgres-user", "baseline", "name for the local postgres user")
	startBaselineStackCmd.Flags().StringVar(&postgresPassword, "postgres-password", "prvdp455", "password for the local postgres user")

	startBaselineStackCmd.Flags().IntVar(&redisContainerPort, "redis-container-port", defaultRedisContainerPort, "container port on which to expose the redis service")
	startBaselineStackCmd.Flags().StringVar(&redisHostname, "redis-hostname", fmt.Sprintf("%s-redis", name), "hostname for the local BPI redis container")
	startBaselineStackCmd.Flags().IntVar(&redisPort, "redis-port", 6379, "host port on which to expose the local redis service")
	startBaselineStackCmd.Flags().StringVar(&redisHosts, "redis-hosts", fmt.Sprintf("%s:%d", redisHostname, redisContainerPort), "list of clustered redis hosts in the local BPI stack")

	startBaselineStackCmd.Flags().BoolVar(&autoRemove, "autoremove", false, "when true, containers are automatically pruned upon exit")
	startBaselineStackCmd.Flags().BoolVar(&prune, "prune", false, "when true, previously-created docker resources are pruned prior to stack initialization")

	startBaselineStackCmd.Flags().StringVar(&logLevel, "log-level", "DEBUG", "log level to set within the running local BPI stack")
	startBaselineStackCmd.Flags().StringVar(&syslogEndpoint, "syslog-endpoint", "", "syslog endpoint to which syslog udp packets will be sent")
	startBaselineStackCmd.Flags().StringVar(&databaseLogging, "database-logging", "false", "when true, query logging is enabled within the local stack")

	startBaselineStackCmd.Flags().StringVar(&jwtSignerPublicKey, "jwt-signer-public-key", "", "PEM-encoded public key of the authorized JWT signer for verifying inbound connection attempts")

	startBaselineStackCmd.Flags().StringVar(&identAPIHost, "ident-host", "ident.provide.services", "hostname of the ident service")
	startBaselineStackCmd.Flags().StringVar(&identAPIScheme, "ident-scheme", "https", "protocol scheme of the ident service")
	startBaselineStackCmd.Flags().IntVar(&identContainerPort, "ident-container-port", 8080, "container port on which to expose the ident service")

	startBaselineStackCmd.Flags().StringVar(&nchainAPIHost, "nchain-host", "nchain.provide.services", "hostname of the nchain service")
	startBaselineStackCmd.Flags().StringVar(&nchainAPIScheme, "nchain-scheme", "https", "protocol scheme of the nchain service")
	startBaselineStackCmd.Flags().IntVar(&nchainContainerPort, "nchain-container-port", 8080, "container port on which to expose the nchain service")

	startBaselineStackCmd.Flags().StringVar(&privacyAPIHost, "privacy-host", "privacy.provide.services", "hostname of the privacy service")
	startBaselineStackCmd.Flags().StringVar(&privacyAPIScheme, "privacy-scheme", "https", "protocol scheme of the privacy service")
	startBaselineStackCmd.Flags().IntVar(&privacyContainerPort, "privacy-container-port", 8080, "container port on which to expose the privacy service")

	startBaselineStackCmd.Flags().StringVar(&vaultAPIHost, "vault-host", "vault.provide.services", "hostname of the vault service")
	startBaselineStackCmd.Flags().StringVar(&vaultAPIScheme, "vault-scheme", "https", "protocol scheme of the vault service")
	startBaselineStackCmd.Flags().IntVar(&vaultContainerPort, "vault-container-port", 8080, "container port on which to expose the vault service")
	startBaselineStackCmd.Flags().StringVar(&vaultRefreshToken, "vault-refresh-token", os.Getenv("VAULT_REFRESH_TOKEN"), "refresh token to vend access tokens for use with vault")
	startBaselineStackCmd.Flags().StringVar(&vaultSealUnsealKey, "vault-seal-unseal-key", os.Getenv("VAULT_SEAL_UNSEAL_KEY"), "seal/unseal key for the vault service")

	startBaselineStackCmd.Flags().BoolVar(&withLocalIdent, "with-local-ident", false, "when true, ident service is run locally")
	startBaselineStackCmd.Flags().IntVar(&identPort, "ident-local-port", 8081, "port for the local ident service")

	startBaselineStackCmd.Flags().BoolVar(&withLocalNChain, "with-local-nchain", false, "when true, nchain service is run locally")
	startBaselineStackCmd.Flags().IntVar(&nchainPort, "nchain-local-port", 8082, "port for the local nchain service")

	startBaselineStackCmd.Flags().BoolVar(&withLocalPrivacy, "with-local-privacy", false, "when true, privacy service is run locally")
	startBaselineStackCmd.Flags().IntVar(&privacyPort, "privacy-local-port", 8083, "port for the local privacy service")

	startBaselineStackCmd.Flags().BoolVar(&withLocalVault, "with-local-vault", false, "when true, vault service is run locally")
	startBaselineStackCmd.Flags().IntVar(&vaultPort, "vault-local-port", 8084, "port for the local vault service")

	startBaselineStackCmd.Flags().StringVar(&common.UserAccessToken, "access-token", "", "access token for the user; useful for running non-interactively")
	startBaselineStackCmd.Flags().StringVar(&organizationRefreshToken, "organization-refresh-token", os.Getenv("PROVIDE_ORGANIZATION_REFRESH_TOKEN"), "refresh token to vend access tokens for use with the local organization")

	startBaselineStackCmd.Flags().BoolVar(&withoutRequireOrganizationKeys, "without-require-organization-keys", false, "when true, no initial keys are required for the organization upon starting the stack")
	startBaselineStackCmd.Flags().BoolVar(&withoutPrivateSubnet, "without-private-subnet", false, "when false, no private subnet will be initialized and services will be attached to appropriate default network if possible")
	startBaselineStackCmd.Flags().BoolVar(&withoutRequireSubjectAccount, "without-require-subject-account", false, "when true, no initial subject account is required upon starting the stack")
	startBaselineStackCmd.Flags().BoolVar(&withoutRequireWorkgroup, "without-require-workgroup", false, "when true, no workgroup is required to start the stack")

	startBaselineStackCmd.Flags().BoolVar(&withLocalBaselineBuild, "with-local-baseline-build", false, "when true, the baseline docker image will be built from your local baseline directory")
	startBaselineStackCmd.Flags().StringVar(&baselineDirPath, "baseline-dir-path", "~/provide/baseline", "the local file path to your baseline directory")

	defaultBaselineOrganizationAddress := "0x"
	if os.Getenv("BASELINE_ORGANIZATION_ADDRESS") != "" {
		defaultBaselineOrganizationAddress = os.Getenv("BASELINE_ORGANIZATION_ADDRESS")
	}

	defaultBaselineRegistryContractAddress := "0x"
	if os.Getenv("BASELINE_REGISTRY_CONTRACT_ADDRESS") != "" {
		defaultBaselineRegistryContractAddress = os.Getenv("BASELINE_REGISTRY_CONTRACT_ADDRESS")
	}

	// defaultNChainBaselineNetworkID := "66d44f30-9092-4182-a3c4-bc02736d6ae5"
	// if os.Getenv("NCHAIN_BASELINE_NETWORK_ID") != "" {
	// 	defaultNChainBaselineNetworkID = os.Getenv("NCHAIN_BASELINE_NETWORK_ID")
	// }

	startBaselineStackCmd.Flags().StringVar(&baselineOrganizationAddress, "organization-address", defaultBaselineOrganizationAddress, "public baseline regsitry address of the organization")
	startBaselineStackCmd.Flags().StringVar(&baselineRegistryContractAddress, "registry-contract-address", defaultBaselineRegistryContractAddress, "public baseline regsitry contract address")
	startBaselineStackCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "baseline workgroup identifier")

	startBaselineStackCmd.Flags().StringVar(&nchainBaselineNetworkID, "nchain-network-id", "", "nchain network id of the baseline mainnet")
	startBaselineStackCmd.Flags().BoolVarP(&Optional, "prompt-all", "", false, "when true, prompts for all optional flags")

	initSORFlags()
}

func initSORFlags() {
	startBaselineStackCmd.Flags().StringVar(&azureServiceBusConnectionString, "azure-servicebus-connection-string", "", "azure service bus connection string")

	startBaselineStackCmd.Flags().StringVar(&salesforceAPIHost, "salesforce-api-host", "", "hostname of the Salesforce API service")
	startBaselineStackCmd.Flags().StringVar(&salesforceAPIScheme, "salesforce-api-scheme", "https", "protocol scheme of the Salesforce API service")
	startBaselineStackCmd.Flags().StringVar(&salesforceAPIPath, "salesforce-api-path", "", "base path of the Salesforce API service")

	startBaselineStackCmd.Flags().StringVar(&sapAPIHost, "sap-api-host", "", "hostname of the internal SAP API service")
	startBaselineStackCmd.Flags().StringVar(&sapAPIScheme, "sap-api-scheme", "https", "protocol scheme of the internal SAP API service")
	startBaselineStackCmd.Flags().StringVar(&sapAPIPath, "sap-api-path", "ubc", "base path of the SAP API service")
	startBaselineStackCmd.Flags().StringVar(&sapAPIUsername, "sap-api-username", "", "username to use for basic authorization against the SAP API service")
	startBaselineStackCmd.Flags().StringVar(&sapAPIPassword, "sap-api-password", "", "password to use for basic authorization against the SAP API service")

	startBaselineStackCmd.Flags().StringVar(&serviceNowAPIHost, "servicenow-api-host", "", "hostname of the ServiceNow service")
	startBaselineStackCmd.Flags().StringVar(&serviceNowAPIScheme, "servicenow-api-scheme", "https", "protocol scheme of the ServiceNow service")
	startBaselineStackCmd.Flags().StringVar(&serviceNowAPIPath, "servicenow-api-path", "api/now/table", "base path of the ServiceNow API")
	startBaselineStackCmd.Flags().StringVar(&serviceNowAPIUsername, "servicenow-api-username", "", "username to use for basic authorization against the ServiceNow API")
	startBaselineStackCmd.Flags().StringVar(&serviceNowAPIPassword, "servicenow-api-password", "", "password to use for basic authorization against the ServiceNow API")
}
