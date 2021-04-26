package stack

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/api/nchain"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const baselineContainerImage = "provide/baseline"
const natsContainerImage = "provide/nats-server"
const natsStreamingContainerImage = "provide/nats-streaming"
const redisContainerImage = "redis"

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
const natsStreamingContainerPort = 4222
const redisContainerPort = 6379

type portMapping struct {
	hostPort      int
	containerPort int
}

var dockerNetworkID string

var name string
var port int
var natsPort int
var natsWebsocketPort int
var natsStreamingPort int
var redisPort int

var apiHostname string
var consumerHostname string
var natsHostname string
var natsStreamingHostname string
var redisHostname string
var redisHosts string

var autoRemove bool
var logLevel string

var baselineOrganizationAddress string
var baselineOrganizationAPIEndpoint string
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

var vaultAPIHost string
var vaultAPIScheme string
var vaultRefreshToken string
var vaultSealUnsealKey string

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

var runBaselineStackCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the baseline stack",
	Long:  `Start a local baseline stack instance and connect to internal systems of record`,
	Run:   runProxy,
}

func runProxy(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	authorizeContext()
	purgeContainers(docker)

	for _, image := range []string{
		baselineContainerImage,
		natsContainerImage,
		natsStreamingContainerImage,
		redisContainerImage,
	} {
		err := pullImage(docker, image)
		if err != nil {
			log.Printf("failed to pull local baseline container image: %s; %s", image, err.Error())
			os.Exit(1)
		}
	}

	configureNetwork(docker)
	common.RequireOrganizationEndpoints(
		func() {
			// run local deps
			runNATS(docker)
			runNATSStreaming(docker)
			runRedis(docker)

			// run proxy
			runProxyAPI(docker)
			runProxyConsumer(docker)

			log.Printf("%s local baseline instance started", name)
		},
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
	log.Printf("authorizing workgroup context")
	authorizeWorkgroupContext()

	log.Printf("authorizing organization context")
	common.AuthorizeOrganizationContext()

	if organizationRefreshToken == "" {
		refreshTokenKey := common.BuildConfigKeyWithOrg(common.APIRefreshTokenConfigKeyPartial, common.OrganizationID)
		if viper.IsSet(refreshTokenKey) {
			// log.Printf("using cached API refresh token for organization: %s\n", common.OrganizationID)
			organizationRefreshToken = viper.GetString(refreshTokenKey)
			if vaultRefreshToken == "" {
				vaultRefreshToken = organizationRefreshToken
			}
		} else {
			log.Printf("failed to resolve refresh token for organization: %s\n", common.OrganizationID)
			os.Exit(1)
		}
	}
}

func authorizeWorkgroupContext() {
	token, err := ident.CreateToken(common.RequireUserAuthToken(), map[string]interface{}{
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

	_, err = ident.GetApplicationDetails(workgroupAccessToken, baselineWorkgroupID, map[string]interface{}{})
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
	}

	orgRegistryContract := contracts[0]
	if orgRegistryContract.Address == nil || *orgRegistryContract.Address == "0x" {
		log.Printf("failed to resolve global organization registry contract; %s", err.Error())
		os.Exit(1)
	}
	baselineRegistryContractAddress = *contracts[0].Address
}

func containerEnvironmentFactory() []string {
	env := make([]string, 0)
	for _, envvar := range []string{
		fmt.Sprintf("BASELINE_ORGANIZATION_ADDRESS=%s", baselineOrganizationAddress),
		fmt.Sprintf("BASELINE_ORGANIZATION_MESSAGING_ENDPOINT=%s", common.MessagingEndpoint),
		fmt.Sprintf("BASELINE_ORGANIZATION_PROXY_ENDPOINT=%s", baselineOrganizationAPIEndpoint),
		fmt.Sprintf("BASELINE_REGISTRY_CONTRACT_ADDRESS=%s", baselineRegistryContractAddress),
		fmt.Sprintf("IDENT_API_HOST=%s", identAPIHost),
		fmt.Sprintf("IDENT_API_SCHEME=%s", identAPIScheme),
		fmt.Sprintf("JWT_SIGNER_PUBLIC_KEY=%s", jwtSignerPublicKey),
		fmt.Sprintf("LOG_LEVEL=%s", logLevel),
		fmt.Sprintf("NATS_CLIENT_PREFIX=%s", name),
		fmt.Sprintf("NATS_STREAMING_URL=%s", fmt.Sprintf("nats://%s:%d", natsStreamingHostname, natsStreamingContainerPort)),
		fmt.Sprintf("NATS_TOKEN=%s", natsAuthToken),
		fmt.Sprintf("NATS_URL=%s", fmt.Sprintf("nats://%s:%d", natsHostname, natsPort)),
		fmt.Sprintf("NCHAIN_API_HOST=%s", nchainAPIHost),
		fmt.Sprintf("NCHAIN_API_SCHEME=%s", nchainAPIScheme),
		fmt.Sprintf("NCHAIN_BASELINE_NETWORK_ID=%s", nchainBaselineNetworkID),
		fmt.Sprintf("PRIVACY_API_HOST=%s", privacyAPIHost),
		fmt.Sprintf("PRIVACY_API_SCHEME=%s", privacyAPIScheme),
		fmt.Sprintf("PROVIDE_ORGANIZATION_ID=%s", common.OrganizationID),
		fmt.Sprintf("PROVIDE_ORGANIZATION_REFRESH_TOKEN=%s", organizationRefreshToken),
		fmt.Sprintf("PROVIDE_SOR_IDENTIFIER=%s", sorID),
		fmt.Sprintf("PROVIDE_SOR_URL=%s", sorURL),
		fmt.Sprintf("PRIVACY_API_SCHEME=%s", privacyAPIScheme),
		fmt.Sprintf("REDIS_HOSTS=%s", redisHosts),
		fmt.Sprintf("VAULT_API_HOST=%s", vaultAPIHost),
		fmt.Sprintf("VAULT_API_SCHEME=%s", vaultAPIScheme),
		fmt.Sprintf("VAULT_REFRESH_TOKEN=%s", vaultRefreshToken),
		fmt.Sprintf("VAULT_SEAL_UNSEAL_KEY=%s", vaultSealUnsealKey),
	} {
		env = append(env, envvar)
	}

	if sapAPIHost != "" && sapAPIUsername != "" && sapAPIPassword != "" {
		for _, envvar := range []string{
			fmt.Sprintf("SAP_API_HOST=%s", sapAPIHost),
			fmt.Sprintf("SAP_API_SCHEME=%s", sapAPIScheme),
			fmt.Sprintf("SAP_API_PATH=%s", sapAPIPath),
			fmt.Sprintf("SAP_API_USERNAE=%s", sapAPIUsername),
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
			fmt.Sprintf("SERVICENOW_API_USERNAE=%s", serviceNowAPIUsername),
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

func runProxyAPI(docker *client.Client) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-api", strings.ReplaceAll(name, " ", "")),
		apiHostname,
		baselineContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", apiHostname, apiContainerPort)},
		[]portMapping{{
			hostPort:      port,
			containerPort: apiContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local baseline API container; %s", err.Error())
		os.Exit(1)
	}
}

func runProxyConsumer(docker *client.Client) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-consumer", strings.ReplaceAll(name, " ", "")),
		consumerHostname,
		baselineContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", apiHostname, port)},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create local baseline consumer container; %s", err.Error())
		os.Exit(1)
	}
}

func runNATS(docker *client.Client) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-nats", strings.ReplaceAll(name, " ", "")),
		natsHostname,
		natsContainerImage,
		nil,
		&[]string{"-auth", natsAuthToken, "-p", fmt.Sprintf("%d", natsPort), "-D", "-V"},
		&[]string{"CMD", "/usr/local/bin/await_tcp.sh", fmt.Sprintf("localhost:%d", natsContainerPort)},
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
}

func runNATSStreaming(docker *client.Client) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-nats-streaming", strings.ReplaceAll(name, " ", "")),
		natsStreamingHostname,
		natsStreamingContainerImage,
		nil,
		&[]string{"-cid", defaultNATSStreamingClusterID, "--auth", natsAuthToken, "-SDV"},
		&[]string{"CMD", "/usr/local/bin/await_tcp.sh", fmt.Sprintf("localhost:%d", natsStreamingContainerPort)},
		[]portMapping{
			{
				hostPort:      natsStreamingPort,
				containerPort: natsStreamingContainerPort,
			},
		}...,
	)

	if err != nil {
		log.Printf("failed to create local baseline NATS streaming container; %s", err.Error())
		os.Exit(1)
	}
}

func runRedis(docker *client.Client) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-redis", strings.ReplaceAll(name, " ", "")),
		redisHostname,
		redisContainerImage,
		nil,
		nil,
		&[]string{"CMD", "redis-cli", "ping"},
		[]portMapping{{
			hostPort:      redisPort,
			containerPort: redisContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local baseline redis container; %s", err.Error())
		os.Exit(1)
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
	entrypoint, cmd, healthcheck *[]string,
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

	containerConfig := &container.Config{
		Env:      containerEnvironmentFactory(),
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

	container, err := docker.ContainerCreate(
		context.Background(),
		containerConfig,
		&container.HostConfig{
			AutoRemove:   autoRemove,
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

func listContainers(docker *client.Client) []types.Container {
	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs([]filters.KeyValuePair{
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-api", strings.ReplaceAll(name, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-consumer", strings.ReplaceAll(name, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-nats", strings.ReplaceAll(name, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-nats-streaming", strings.ReplaceAll(name, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-redis", strings.ReplaceAll(name, " ", "")),
			},
		}...),
	})
	if err != nil {
		log.Printf("failed to list containers; %s", err.Error())
		os.Exit(1)
	}

	return containers
}

func init() {
	runBaselineStackCmd.Flags().StringVar(&name, "name", "baseline-local", "name of the baseline stack instance")

	runBaselineStackCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	runBaselineStackCmd.MarkFlagRequired("organization")

	runBaselineStackCmd.Flags().StringVar(&baselineOrganizationAPIEndpoint, "api-endpoint", "", "local baseline API endpoint for use by one or more authorized systems of record")
	runBaselineStackCmd.Flags().StringVar(&common.MessagingEndpoint, "messaging-endpoint", "", "public messaging endpoint used for sending and receiving protocol messages")
	runBaselineStackCmd.Flags().BoolVar(&common.Tunnel, "tunnel", false, "when true, a tunnel is established to expose the API and messaging endpoints to the WAN")
	runBaselineStackCmd.Flags().BoolVar(&common.ExposeAPITunnel, "api-tunnel", false, "when true, a tunnel is established to expose the API endpoint to the WAN")
	runBaselineStackCmd.Flags().BoolVar(&common.ExposeMessagingTunnel, "messaging-tunnel", false, "when true, a tunnel is established to expose the messaging endpoint to the WAN")

	runBaselineStackCmd.Flags().StringVar(&sorID, "sor", "", "primary internal system of record identifier being baselined")
	runBaselineStackCmd.Flags().StringVar(&sorURL, "sor-url", "https://", "url of the primary internal system of record being baselined")

	runBaselineStackCmd.Flags().StringVar(&apiHostname, "hostname", fmt.Sprintf("%s-api", name), "hostname for the local baseline API container")
	runBaselineStackCmd.Flags().IntVar(&port, "port", 8080, "port to expose on the local baseline API container")

	runBaselineStackCmd.Flags().StringVar(&consumerHostname, "consumer-hostname", fmt.Sprintf("%s-consumer", name), "hostname for the local baseline consumer container")
	runBaselineStackCmd.Flags().StringVar(&natsHostname, "nats-hostname", fmt.Sprintf("%s-nats", name), "hostname for the local baseline NATS container")
	runBaselineStackCmd.Flags().IntVar(&natsPort, "nats-port", 4222, "port to expose on the local baseline NATS container")
	runBaselineStackCmd.Flags().IntVar(&natsWebsocketPort, "nats-ws-port", 4221, "websocket port to expose on the local baseline NATS container")
	runBaselineStackCmd.Flags().StringVar(&natsAuthToken, "nats-auth-token", "testtoken", "authorization token for the local baseline NATS service; will be passed as the -auth argument to NATS")

	runBaselineStackCmd.Flags().StringVar(&natsStreamingHostname, "nats-streaming-hostname", fmt.Sprintf("%s-nats-streaming", name), "hostname for the local baseline NATS streaming container")
	runBaselineStackCmd.Flags().IntVar(&natsStreamingPort, "nats-streaming-port", 4220, "port to expose on the local baseline NATS streaming container")

	runBaselineStackCmd.Flags().StringVar(&redisHostname, "redis-hostname", fmt.Sprintf("%s-redis", name), "hostname for the local baseline redis container")
	runBaselineStackCmd.Flags().IntVar(&redisPort, "redis-port", 6379, "port to expose on the local baseline redis container")
	runBaselineStackCmd.Flags().StringVar(&redisHosts, "redis-hosts", fmt.Sprintf("%s:%d", redisHostname, redisPort), "list of clustered redis hosts in the local baseline stack")

	runBaselineStackCmd.Flags().BoolVar(&autoRemove, "autoremove", false, "when true, containers are automatically pruned upon exit")
	runBaselineStackCmd.Flags().StringVar(&logLevel, "log-level", "DEBUG", "log level to set within the running local baseline stack")

	runBaselineStackCmd.Flags().StringVar(&jwtSignerPublicKey, "jwt-signer-public-key", "", "PEM-encoded public key of the authorized JWT signer for verifying inbound connection attempts")

	runBaselineStackCmd.Flags().StringVar(&identAPIHost, "ident-host", "ident.provide.services", "hostname of the ident service")
	runBaselineStackCmd.Flags().StringVar(&identAPIScheme, "ident-scheme", "https", "protocol scheme of the ident service")

	runBaselineStackCmd.Flags().StringVar(&nchainAPIHost, "nchain-host", "nchain.provide.services", "hostname of the nchain service")
	runBaselineStackCmd.Flags().StringVar(&nchainAPIScheme, "nchain-scheme", "https", "protocol scheme of the nchain service")

	runBaselineStackCmd.Flags().StringVar(&privacyAPIHost, "privacy-host", "privacy.provide.services", "hostname of the privacy service")
	runBaselineStackCmd.Flags().StringVar(&privacyAPIScheme, "privacy-scheme", "https", "protocol scheme of the privacy service")

	runBaselineStackCmd.Flags().StringVar(&vaultAPIHost, "vault-host", "vault.provide.services", "hostname of the vault service")
	runBaselineStackCmd.Flags().StringVar(&vaultAPIScheme, "vault-scheme", "https", "protocol scheme of the vault service")
	runBaselineStackCmd.Flags().StringVar(&vaultRefreshToken, "vault-refresh-token", os.Getenv("VAULT_REFRESH_TOKEN"), "refresh token to vend access tokens for use with vault")
	runBaselineStackCmd.Flags().StringVar(&vaultSealUnsealKey, "vault-seal-unseal-key", os.Getenv("VAULT_SEAL_UNSEAL_KEY"), "seal/unseal key for the vault service")

	runBaselineStackCmd.Flags().StringVar(&organizationRefreshToken, "organization-refresh-token", os.Getenv("PROVIDE_ORGANIZATION_REFRESH_TOKEN"), "refresh token to vend access tokens for use with the local organization")

	defaultBaselineOrganizationAddress := "0x"
	if os.Getenv("BASELINE_ORGANIZATION_ADDRESS") != "" {
		defaultBaselineOrganizationAddress = os.Getenv("BASELINE_ORGANIZATION_ADDRESS")
	}

	defaultBaselineRegistryContractAddress := "0x"
	if os.Getenv("BASELINE_REGISTRY_CONTRACT_ADDRESS") != "" {
		defaultBaselineRegistryContractAddress = os.Getenv("BASELINE_REGISTRY_CONTRACT_ADDRESS")
	}

	defaultNChainBaselineNetworkID := "66d44f30-9092-4182-a3c4-bc02736d6ae5"
	if os.Getenv("NCHAIN_BASELINE_NETWORK_ID") != "" {
		defaultNChainBaselineNetworkID = os.Getenv("NCHAIN_BASELINE_NETWORK_ID")
	}

	runBaselineStackCmd.Flags().StringVar(&baselineOrganizationAddress, "organization-address", defaultBaselineOrganizationAddress, "public baseline regsitry address of the organization")
	runBaselineStackCmd.Flags().StringVar(&baselineRegistryContractAddress, "registry-contract-address", defaultBaselineRegistryContractAddress, "public baseline regsitry contract address")
	runBaselineStackCmd.Flags().StringVar(&baselineWorkgroupID, "workgroup", "", "baseline workgroup identifier")

	runBaselineStackCmd.Flags().StringVar(&nchainBaselineNetworkID, "nchain-network-id", defaultNChainBaselineNetworkID, "nchain network id of the baseline mainnet")

	initSORFlags()
}

func initSORFlags() {
	// runBaselineStackCmd.Flags().StringVar(&salesforceAPIHost, "salesforce-api-host", "", "hostname of the Salesforce API service")
	// runBaselineStackCmd.Flags().StringVar(&salesforceAPIScheme, "salesforce-api-scheme", "https", "protocol scheme of the Salesforce API service")
	// runBaselineStackCmd.Flags().StringVar(&salesforceAPIPath, "salesforce-api-path", "", "base path of the Salesforce API service")

	runBaselineStackCmd.Flags().StringVar(&sapAPIHost, "sap-api-host", "", "hostname of the internal SAP API service")
	runBaselineStackCmd.Flags().StringVar(&sapAPIScheme, "sap-api-scheme", "https", "protocol scheme of the internal SAP API service")
	runBaselineStackCmd.Flags().StringVar(&sapAPIPath, "sap-api-path", "ubc", "base path of the SAP API service")
	runBaselineStackCmd.Flags().StringVar(&sapAPIUsername, "sap-api-username", "", "username to use for basic authorization against the SAP API service")
	runBaselineStackCmd.Flags().StringVar(&sapAPIPassword, "sap-api-password", "", "password to use for basic authorization against the SAP API service")

	runBaselineStackCmd.Flags().StringVar(&serviceNowAPIHost, "servicenow-api-host", "", "hostname of the ServiceNow service")
	runBaselineStackCmd.Flags().StringVar(&serviceNowAPIScheme, "servicenow-api-scheme", "https", "protocol scheme of the ServiceNow service")
	runBaselineStackCmd.Flags().StringVar(&serviceNowAPIPath, "servicenow-api-path", "api/now/table", "base path of the ServiceNow API")
	runBaselineStackCmd.Flags().StringVar(&serviceNowAPIUsername, "servicenow-api-username", "", "username to use for basic authorization against the ServiceNow API")
	runBaselineStackCmd.Flags().StringVar(&serviceNowAPIPassword, "servicenow-api-password", "", "password to use for basic authorization against the ServiceNow API")
}
