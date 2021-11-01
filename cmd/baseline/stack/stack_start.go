package stack

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
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
	"github.com/docker/go-connections/nat"
	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/cmd/common"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/nchain"
	"github.com/provideplatform/provide-go/api/vault"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const baselineContainerImage = "provide/baseline"
const identContainerImage = "provide/ident"
const nchainContainerImage = "provide/nchain"
const privacyContainerImage = "provide/privacy"
const vaultContainerImage = "provide/vault"
const postgresContainerImage = "postgres"
const natsContainerImage = "provide/nats-server:2.5.0-PRVD"
const redisContainerImage = "redis"
const defaultNatsServerName = "prvd"
const defaultNatsReachabilityTimeout = time.Second * 5
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

const apiContainerPort = 8080
const natsContainerPort = 4222
const natsWebsocketContainerPort = 4221
const postgresContainerPort = 5432
const redisContainerPort = 6379

type portMapping struct {
	hostPort      int
	containerPort int
}

var dockerNetworkID string
var Optional bool
var name string
var port int
var identPort int
var nchainPort int
var privacyPort int
var vaultPort int
var natsPort int
var natsWebsocketPort int
var natsWebsocketTLS bool
var postgresPort int
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
var postgresHostname string
var redisHostname string
var redisHosts string

var autoRemove bool

var logLevel string
var syslogEndpoint string

var baselineOrganizationAddress string

// var baselineOrganizationAPIEndpoint string
var baselineRegistryContractAddress string
var baselineWorkgroupID string

var nchainBaselineNetworkID string

var jwtSignerPublicKey string
var natsAuthToken string

var identAPIHost string
var identAPIScheme string

var nchainAPIHost string
var nchainAPIScheme string

var workgroupAccessToken string
var organizationRefreshToken string

var privacyAPIHost string
var privacyAPIScheme string

var sorID string
var sorURL string
var sorOrganizationCode string

var vaultAPIHost string
var vaultAPIScheme string
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

var startBaselineStackCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the baseline stack",
	Long:  `Start a local baseline stack instance and connect to internal systems of record`,
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

	go common.PurgeContainers(docker, name)

	authorizeContext()
	log.Printf("context authorized")
	sorPrompt()
	log.Printf("sorID: %s", sorID)
	tunnelAPIPrompt()
	log.Printf("expose API tunnel: %t", common.ExposeAPITunnel)
	tunnelMessagingPrompt()
	log.Printf("expose messaging tunnel: %t", common.ExposeMessagingTunnel)

	wg := &sync.WaitGroup{}

	images := make([]string, 0)
	images = append(
		images,
		baselineContainerImage,
		natsContainerImage,
		redisContainerImage,
	)

	log.Printf("withLocalIdent: %t, withLocalNChain: %t, withLocalPrivacy: %t, withLocalVault: %t", withLocalIdent, withLocalNChain, withLocalPrivacy, withLocalVault)
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
			log.Printf("image eq: '%s'", img)
			if img != "provide/baseline" {
				err := pullImage(docker, img)
				if err != nil {
						log.Printf("failed to pull local baseline container image: %s; %s", img, err.Error())
						os.Exit(1)
				}
			} else {
				log.Printf("skipping %s pull from dockerhub", img)
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
			go runNATS(docker, wg)

			// FIXME-- DRY this up...
			natsReachable := false
			for !natsReachable {
				host := fmt.Sprintf("localhost:%v", natsPort)
				log.Printf("%s", host)
				conn, err := net.DialTimeout("tcp", host, defaultNatsReachabilityTimeout)
				if err == nil {
					conn.Close()
					natsReachable = true
				}
			}

			wg.Add(1)
			go runRedis(docker, wg)

			// FIXME-- DRY this up...
			redisReachable := false
			for !redisReachable {
				host := fmt.Sprintf("localhost:%v", redisPort)
				log.Printf("%s", host)
				conn, err := net.DialTimeout("tcp", host, defaultRedisReachabilityTimeout)
				if err == nil {
					conn.Close()
					redisReachable = true
				}
			}

			if withLocalIdent || withLocalNChain || withLocalPrivacy || withLocalVault {
				wg.Add(1)
				go runPostgres(docker, wg)
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

			// run proxy
			wg.Add(1)
			go runProxyAPI(docker, wg)

			wg.Add(1)
			go runProxyConsumer(docker, wg)

			wg.Wait()
			log.Printf("%s local baseline instance started", name)
		},
		func(reason *string) {
			if reason != nil {
				log.Printf(*reason)
				common.PurgeContainers(docker, name)
				common.PurgeNetwork(docker, name)
			}
		},
		port,
		natsPort,
	)
}

func configureNetwork(docker *client.Client) {
	network, err := docker.NetworkCreate(
		context.Background(),
		name,
		types.NetworkCreate{
			// CheckDuplicate bool
			Driver: "bridge",
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
		},
	)

	if err != nil {
		log.Printf("failed to setup docker network; %s", err.Error())
		os.Exit(1)
	}

	dockerNetworkID = network.ID
	log.Printf("configured network for local baseline instance: %s", name)
}

func authorizeContext() {
	// log.Printf("authorizing workgroup context")
	authorizeWorkgroupContext()

	// log.Printf("authorizing organization context")
	common.AuthorizeOrganizationContext(false)

	if organizationRefreshToken == "" {
		refreshTokenKey := common.BuildConfigKeyWithOrg(common.APIRefreshTokenConfigKeyPartial, common.OrganizationID)
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
	}
}

func authorizeWorkgroupContext() {
	if baselineWorkgroupID == "" {
		err := common.RequireWorkgroup()
		if err != nil {
			log.Printf("failed to require workgroup; %s", err.Error())
			os.Exit(1)
		}
		baselineWorkgroupID = common.ApplicationID
	}

	token, err := ident.CreateToken(common.RequireUserAccessToken(), map[string]interface{}{
		"scope":          "offline_access",
		"application_id": baselineWorkgroupID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of workgroup %s; %s", baselineWorkgroupID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		workgroupAccessToken = *token.AccessToken
	}

	workgroup, err := ident.GetApplicationDetails(workgroupAccessToken, baselineWorkgroupID, map[string]interface{}{})
	if err != nil {
		log.Printf("failed to resolve workgroup: %s; %s", baselineWorkgroupID, err.Error())
		os.Exit(1)
	}

	contracts, err := nchain.ListContracts(workgroupAccessToken, map[string]interface{}{
		"type": "organization-registry",
	})
	if err != nil {
		log.Printf("failed to resolve global organization registry contract; %s", err.Error())
		os.Exit(1)
	} else if len(contracts) == 0 {
		log.Printf("failed to resolve global organization registry contract")
		os.Exit(1)
	}

	if nchainBaselineNetworkID == "" {
		if workgroup.NetworkID != uuid.Nil {
			nchainBaselineNetworkID = workgroup.NetworkID.String()
		} else {
			err := common.RequirePublicNetwork()
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
		fmt.Sprintf("BASELINE_ORGANIZATION_ADDRESS=%s", baselineOrganizationAddress),
		fmt.Sprintf("BASELINE_ORGANIZATION_MESSAGING_ENDPOINT=%s", common.MessagingEndpoint),
		fmt.Sprintf("BASELINE_ORGANIZATION_PROXY_ENDPOINT=%s", common.APIEndpoint),
		fmt.Sprintf("BASELINE_REGISTRY_CONTRACT_ADDRESS=%s", baselineRegistryContractAddress),
		fmt.Sprintf("BASELINE_WORKGROUP_ID=%s", baselineWorkgroupID),
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
		fmt.Sprintf("PROVIDE_ORGANIZATION_ID=%s", common.OrganizationID),
		fmt.Sprintf("PROVIDE_ORGANIZATION_REFRESH_TOKEN=%s", organizationRefreshToken),
		fmt.Sprintf("PROVIDE_SOR_IDENTIFIER=%s", sorID),
		fmt.Sprintf("PROVIDE_SOR_ORGANIZATION_CODE=%s", sorOrganizationCode),
		fmt.Sprintf("PROVIDE_SOR_URL=%s", sorURL),
		fmt.Sprintf("PRIVACY_API_SCHEME=%s", privacyAPIScheme),
		fmt.Sprintf("REDIS_HOSTS=%s", redisHosts),
		fmt.Sprintf("SYSLOG_ENDPOINT=%s", syslogEndpoint),
		fmt.Sprintf("VAULT_API_HOST=%s", vaultAPIHost),
		fmt.Sprintf("VAULT_API_SCHEME=%s", vaultAPIScheme),
		fmt.Sprintf("VAULT_REFRESH_TOKEN=%s", vaultRefreshToken),
		fmt.Sprintf("VAULT_SEAL_UNSEAL_KEY=%s", vaultSealUnsealKey),
	} {
		env = append(env, envvar)
	}

	if listenPort != nil {
		env = append(env, fmt.Sprintf("PORT=%d", *listenPort))
	}

	if azureServiceBusConnectionString != "" {
		for _, envvar := range []string{
			fmt.Sprintf("AZURE_SERVICE_BUS_CONNECTION_STRING=%s", azureServiceBusConnectionString),
		} {
			env = append(env, envvar)
		}
	}

	if sapAPIHost != "" && sapAPIUsername != "" && sapAPIPassword != "" {
		for _, envvar := range []string{
			fmt.Sprintf("SAP_API_HOST=%s", sapAPIHost),
			fmt.Sprintf("SAP_API_SCHEME=%s", sapAPIScheme),
			fmt.Sprintf("SAP_API_PATH=%s", sapAPIPath),
			fmt.Sprintf("SAP_API_USERNAME=%s", sapAPIUsername),
			fmt.Sprintf("SAP_API_PASSWORD=%s", sapAPIPassword),
		} {
			env = append(env, envvar)
		}
	} else if sorID == "sap" && sorURL != "" {
		_url, err := url.Parse(sorURL)
		if err != nil {
			log.Printf("WARNING: system of record url invalid; %s", err.Error())
		}
		for _, envvar := range []string{
			fmt.Sprintf("SAP_API_HOST=%s", _url.Host),
			fmt.Sprintf("SAP_API_SCHEME=%s", _url.Scheme),
		} {
			env = append(env, envvar)
		}

		if _url.Path != "" {
			env = append(env, fmt.Sprintf("SAP_API_PATH=%s", strings.TrimLeft(_url.Path, "/")))
		}
	}

	if serviceNowAPIHost != "" && serviceNowAPIUsername != "" && serviceNowAPIPassword != "" {
		for _, envvar := range []string{
			fmt.Sprintf("SERVICENOW_API_HOST=%s", serviceNowAPIHost),
			fmt.Sprintf("SERVICENOW_API_SCHEME=%s", serviceNowAPIScheme),
			fmt.Sprintf("SERVICENOW_API_PATH=%s", serviceNowAPIPath),
			fmt.Sprintf("SERVICENOW_API_USERNAME=%s", serviceNowAPIUsername),
			fmt.Sprintf("SERVICENOW_API_PASSWORD=%s", serviceNowAPIPassword),
		} {
			env = append(env, envvar)
		}
	} else if sorID == "servicenow" || sorID == "snow" && sorURL != "" {
		_url, err := url.Parse(sorURL)
		if err != nil {
			log.Printf("WARNING: system of record url invalid; %s", err.Error())
		}
		for _, envvar := range []string{
			fmt.Sprintf("SERVICENOW_API_HOST=%s", _url.Host),
			fmt.Sprintf("SERVICENOW_API_SCHEME=%s", _url.Scheme),
		} {
			env = append(env, envvar)
		}

		if _url.Path != "" {
			env = append(env, fmt.Sprintf("SERVICENOW_API_PATH=%s", strings.TrimLeft(_url.Path, "/")))
		}
	}

	return env
}

func runProxyAPI(docker *client.Client, wg *sync.WaitGroup) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-api", strings.ReplaceAll(name, " ", "")),
		apiHostname,
		baselineContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", apiHostname, apiContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      port,
			containerPort: apiContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local baseline API container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runProxyConsumer(docker *client.Client, wg *sync.WaitGroup) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-consumer", strings.ReplaceAll(name, " ", "")),
		consumerHostname,
		baselineContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", apiHostname, port)},
		nil,
		map[string]string{},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create local baseline consumer container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runIdentAPI(docker *client.Client, wg *sync.WaitGroup) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-ident-api", strings.ReplaceAll(name, " ", "")),
		identHostname,
		identContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", identHostname, apiContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      identPort,
			containerPort: apiContainerPort,
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
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-ident-consumer", strings.ReplaceAll(name, " ", "")),
		identConsumerHostname,
		identContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", identHostname, apiContainerPort)},
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
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-nchain-api", strings.ReplaceAll(name, " ", "")),
		nchainHostname,
		nchainContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", nchainHostname, apiContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      nchainPort,
			containerPort: apiContainerPort,
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
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-nchain-consumer", strings.ReplaceAll(name, " ", "")),
		nchainConsumerHostname,
		nchainContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", nchainHostname, apiContainerPort)},
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
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-statsdaemon", strings.ReplaceAll(name, " ", "")),
		nchainStatsdaemonHostname,
		nchainContainerImage,
		&[]string{"./ops/run_statsdaemon.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", nchainHostname, apiContainerPort)},
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
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-reachabilitydaemon", strings.ReplaceAll(name, " ", "")),
		nchainReachabilitydaemonHostname,
		nchainContainerImage,
		&[]string{"./ops/run_reachabilitydaemon.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", nchainHostname, apiContainerPort)},
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
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-privacy-api", strings.ReplaceAll(name, " ", "")),
		privacyHostname,
		privacyContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", privacyHostname, apiContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      privacyPort,
			containerPort: apiContainerPort,
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
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-privacy-consumer", strings.ReplaceAll(name, " ", "")),
		privacyConsumerHostname,
		privacyContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", privacyHostname, apiContainerPort)},
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
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-vault-api", strings.ReplaceAll(name, " ", "")),
		vaultHostname,
		vaultContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", vaultHostname, apiContainerPort)},
		nil,
		map[string]string{},
		[]portMapping{{
			hostPort:      vaultPort,
			containerPort: apiContainerPort,
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

func writeNATSConfig() {
	cfg := []byte("max_payload: 100Mb\nmax_pending: 104857600\n")
	if !natsWebsocketTLS {
		cfg = []byte("max_payload: 100Mb\nmax_pending: 104857600\nwebsocket {\n    listen: \"0.0.0.0:4221\"\n    no_tls: true\n}\n")
	}
	tmp := filepath.Join(os.TempDir(), "nats-server.conf")
	err := ioutil.WriteFile(tmp, cfg, 0644)
	if err != nil {
		log.Printf("failed to write local nats-server.conf; %s", err.Error())
		os.Exit(1)
	}
}

func runNATS(docker *client.Client, wg *sync.WaitGroup) {
	writeNATSConfig()

	_, err := runContainer(
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
		&[]string{"CMD", "/usr/local/bin/await_tcp.sh", fmt.Sprintf("localhost:%d", natsContainerPort)},
		nil,
		map[string]string{
			filepath.Join(os.TempDir(), "nats-server.conf"): "/etc/nats-server.conf",
		},
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
		log.Printf("failed to create local baseline NATS container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func runPostgres(docker *client.Client, wg *sync.WaitGroup) {
	_, err := runContainer(
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
	_, err := runContainer(
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
		log.Printf("failed to create local baseline redis container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func pullImage(docker *client.Client, image string) error {
	log.Printf("pulling local baseline container image: %s", image)
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
) (*container.ContainerCreateCreatedBody, error) {
	log.Printf("running local baseline container image: %s", image)
	portBinding := nat.PortMap{}
	for _, mapping := range ports {
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

	container, err := docker.ContainerCreate(
		context.Background(),
		containerConfig,
		&container.HostConfig{
			AutoRemove:   autoRemove,
			ExtraHosts:   []string{"host.docker.internal:host-gateway"},
			Mounts:       mountedVolumes,
			NetworkMode:  "bridge",
			PortBindings: portBinding,
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		},
		&network.NetworkingConfig{},
		strings.ReplaceAll(name, " ", ""),
	)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	err = docker.NetworkConnect(
		context.Background(),
		dockerNetworkID,
		container.ID,
		&network.EndpointSettings{},
	)
	if err != nil {
		return nil, err
	}

	return &container, nil
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
	if common.ExposeAPITunnel || common.APIEndpoint != "" {
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
		common.ExposeAPITunnel = true
	}
}

func tunnelMessagingPrompt() {
	if common.ExposeMessagingTunnel || common.MessagingEndpoint != "" {
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
	startBaselineStackCmd.Flags().StringVar(&name, "name", "baseline-local", "name of the baseline stack instance")

	startBaselineStackCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	// runBaselineStackCmd.MarkFlagRequired("organization")

	startBaselineStackCmd.Flags().StringVar(&common.APIEndpoint, "api-endpoint", "", "local baseline API endpoint for use by one or more authorized systems of record")
	startBaselineStackCmd.Flags().StringVar(&common.MessagingEndpoint, "messaging-endpoint", "", "public messaging endpoint used for sending and receiving protocol messages")
	startBaselineStackCmd.Flags().BoolVar(&common.Tunnel, "tunnel", false, "when true, a tunnel is established to expose the API and messaging endpoints to the WAN")
	startBaselineStackCmd.Flags().BoolVar(&common.ExposeAPITunnel, "api-tunnel", false, "when true, a tunnel is established to expose the API endpoint to the WAN")
	startBaselineStackCmd.Flags().BoolVar(&common.ExposeMessagingTunnel, "messaging-tunnel", false, "when true, a tunnel is established to expose the messaging endpoint to the WAN")

	startBaselineStackCmd.Flags().StringVar(&sorID, "sor", "", "primary internal system of record identifier being baselined")
	startBaselineStackCmd.Flags().StringVar(&sorURL, "sor-url", "https://", "url of the primary internal system of record being baselined")
	startBaselineStackCmd.Flags().StringVar(&sorOrganizationCode, "sor-organization-code", "", "organization code specific to the system of record")

	startBaselineStackCmd.Flags().StringVar(&apiHostname, "hostname", fmt.Sprintf("%s-api", name), "hostname for the local baseline API container")
	startBaselineStackCmd.Flags().IntVar(&port, "port", 8080, "host port on which to expose the local baseline API service")

	startBaselineStackCmd.Flags().StringVar(&consumerHostname, "consumer-hostname", fmt.Sprintf("%s-consumer", name), "hostname for the local baseline consumer container")
	startBaselineStackCmd.Flags().StringVar(&natsHostname, "nats-hostname", fmt.Sprintf("%s-nats", name), "hostname for the local baseline NATS container")
	startBaselineStackCmd.Flags().IntVar(&natsPort, "nats-port", 4222, "host port on which to expose the local NATS service")
	startBaselineStackCmd.Flags().BoolVar(&natsWebsocketTLS, "nats-ws-tls", false, "when true, NATS websocket service uses TLS")
	startBaselineStackCmd.Flags().IntVar(&natsWebsocketPort, "nats-ws-port", 4221, "host port on which to expose the local NATS websocket service")
	startBaselineStackCmd.Flags().StringVar(&natsAuthToken, "nats-auth-token", "testtoken", "authorization token for the local baseline NATS service; will be passed as the -auth argument to NATS")

	startBaselineStackCmd.Flags().StringVar(&postgresHostname, "postgres-hostname", fmt.Sprintf("%s-postgres", name), "hostname for the local postgres container")
	startBaselineStackCmd.Flags().IntVar(&postgresPort, "postgres-port", 5432, "host port on which to expose the local postgres service")

	startBaselineStackCmd.Flags().StringVar(&redisHostname, "redis-hostname", fmt.Sprintf("%s-redis", name), "hostname for the local baseline redis container")
	startBaselineStackCmd.Flags().IntVar(&redisPort, "redis-port", 6379, "host port on which to expose the local redis service")
	startBaselineStackCmd.Flags().StringVar(&redisHosts, "redis-hosts", fmt.Sprintf("%s:%d", redisHostname, redisContainerPort), "list of clustered redis hosts in the local baseline stack")

	startBaselineStackCmd.Flags().BoolVar(&autoRemove, "autoremove", false, "when true, containers are automatically pruned upon exit")

	startBaselineStackCmd.Flags().StringVar(&logLevel, "log-level", "DEBUG", "log level to set within the running local baseline stack")
	startBaselineStackCmd.Flags().StringVar(&syslogEndpoint, "syslog-endpoint", "", "syslog endpoint to which syslog udp packets will be sent")

	startBaselineStackCmd.Flags().StringVar(&jwtSignerPublicKey, "jwt-signer-public-key", "", "PEM-encoded public key of the authorized JWT signer for verifying inbound connection attempts")

	startBaselineStackCmd.Flags().StringVar(&identAPIHost, "ident-host", "ident.provide.services", "hostname of the ident service")
	startBaselineStackCmd.Flags().StringVar(&identAPIScheme, "ident-scheme", "https", "protocol scheme of the ident service")

	startBaselineStackCmd.Flags().StringVar(&nchainAPIHost, "nchain-host", "nchain.provide.services", "hostname of the nchain service")
	startBaselineStackCmd.Flags().StringVar(&nchainAPIScheme, "nchain-scheme", "https", "protocol scheme of the nchain service")

	startBaselineStackCmd.Flags().StringVar(&privacyAPIHost, "privacy-host", "privacy.provide.services", "hostname of the privacy service")
	startBaselineStackCmd.Flags().StringVar(&privacyAPIScheme, "privacy-scheme", "https", "protocol scheme of the privacy service")

	startBaselineStackCmd.Flags().StringVar(&vaultAPIHost, "vault-host", "vault.provide.services", "hostname of the vault service")
	startBaselineStackCmd.Flags().StringVar(&vaultAPIScheme, "vault-scheme", "https", "protocol scheme of the vault service")
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

	startBaselineStackCmd.Flags().StringVar(&organizationRefreshToken, "organization-refresh-token", os.Getenv("PROVIDE_ORGANIZATION_REFRESH_TOKEN"), "refresh token to vend access tokens for use with the local organization")

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
	startBaselineStackCmd.Flags().StringVar(&baselineWorkgroupID, "workgroup", "", "baseline workgroup identifier")

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
