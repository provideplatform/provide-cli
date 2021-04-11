package baseline

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/provideservices/provide-cli/cmd/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const baselineProxyContainerImage = "provide/baseline"
const natsContainerImage = "provide/nats-server"
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

type portMapping struct {
	hostPort      int
	containerPort int
}

var name string
var port int
var natsPort int
var natsWebsocketPort int
var redisPort int

var apiHostname string
var natsHostname string
var redisHostname string
var redisHosts string

var autoRemove bool
var logLevel string

var baselineOrganizationAddress string
var baselineRegistryContractAddress string
var nchainBaselineNetworkID string

var jwtSignerPublicKey string
var natsAuthToken string

var identAPIHost string
var identAPIScheme string

var nchainAPIHost string
var nchainAPIScheme string

var organizationRefreshToken string

var privacyAPIHost string
var privacyAPIScheme string

var sorID string
var sorURL string

var vaultAPIHost string
var vaultAPIScheme string
var vaultRefreshToken string
var vaultSealUnsealKey string

var runBaselineProxyCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the baseline proxy",
	Long:  `Start a local baseline proxy instance and connect to internal systems of record`,
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
		baselineProxyContainerImage,
		natsContainerImage,
		redisContainerImage,
	} {
		err := pullImage(docker, image)
		if err != nil {
			log.Printf("failed to pull proxy container image: %s; %s", image, err.Error())
			os.Exit(1)
		}
	}

	// run local deps
	runNATS(docker)
	runRedis(docker)

	// run proxy
	runProxyAPI(docker)
	runProxyConsumer(docker)

	log.Printf("%s proxy instance started", name)
}

func authorizeContext() {
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

func containerEnvironmentFactory() []string {
	return []string{
		fmt.Sprintf("BASELINE_ORGANIZATION_ADDRESS=%s", baselineOrganizationAddress),
		fmt.Sprintf("BASELINE_REGISTRY_CONTRACT_ADDRESS=%s", baselineRegistryContractAddress),
		fmt.Sprintf("IDENT_API_HOST=%s", identAPIHost),
		fmt.Sprintf("IDENT_API_SCHEME=%s", identAPIScheme),
		fmt.Sprintf("JWT_SIGNER_PUBLIC_KEY=%s", jwtSignerPublicKey),
		fmt.Sprintf("LOG_LEVEL=%s", logLevel),
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
	}
}

func runProxyAPI(docker *client.Client) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-api", strings.ReplaceAll(name, " ", "")),
		apiHostname,
		baselineProxyContainerImage,
		&[]string{"./ops/run_api.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", apiHostname, port)},
		[]portMapping{{
			hostPort:      port,
			containerPort: port,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create baseline proxy API container; %s", err.Error())
		os.Exit(1)
	}
}

func runProxyConsumer(docker *client.Client) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-consumer", strings.ReplaceAll(name, " ", "")),
		natsHostname,
		baselineProxyContainerImage,
		&[]string{"./ops/run_consumer.sh"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", apiHostname, port)},
		[]portMapping{}...,
	)

	if err != nil {
		log.Printf("failed to create baseline proxy consumer container; %s", err.Error())
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
		&[]string{"CMD", "/usr/local/bin/await_tcp.sh", fmt.Sprintf("localhost:%d", natsPort)},
		[]portMapping{
			{
				hostPort:      natsPort,
				containerPort: natsPort,
			},
			{
				hostPort:      natsWebsocketPort,
				containerPort: natsWebsocketPort,
			},
		}...,
	)

	if err != nil {
		log.Printf("failed to create baseline proxy NATS container; %s", err.Error())
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
			containerPort: redisPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create baseline proxy Redis container; %s", err.Error())
		os.Exit(1)
	}
}

func pullImage(docker *client.Client, image string) error {
	reader, err := docker.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, reader)

	return nil
}

func runContainer(
	docker *client.Client,
	name, hostname, image string,
	entrypoint, cmd, healthcheck *[]string,
	ports ...portMapping,
) (*container.ContainerCreateCreatedBody, error) {
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
			PortBindings: portBinding,
		},
		nil,
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
	runBaselineProxyCmd.Flags().StringVar(&name, "name", "baseline-proxy", "name of the baseline proxy instance")
	runBaselineProxyCmd.Flags().IntVar(&port, "port", 8080, "local API port to expose on the proxy")
	runBaselineProxyCmd.Flags().IntVar(&natsPort, "nats-port", 4222, "local NATS port to expose on the proxy")
	runBaselineProxyCmd.Flags().IntVar(&natsWebsocketPort, "nats-ws-port", 4221, "local NATS websocket port to expose on the proxy")
	runBaselineProxyCmd.Flags().IntVar(&redisPort, "redis-port", 6379, "local NATS port to expose on the proxy")

	runBaselineProxyCmd.Flags().StringVar(&apiHostname, "hostname", "baseline-proxy-api", "hostname for the proxy API container")
	runBaselineProxyCmd.Flags().StringVar(&natsHostname, "nats-hostname", "baseline-proxy-nats", "hostname for the proxy NATS container")
	runBaselineProxyCmd.Flags().StringVar(&redisHostname, "redis-hostname", "baseline-proxy-redis", "hostname for the proxy Redis container")
	runBaselineProxyCmd.Flags().StringVar(&redisHosts, "redis-hosts", fmt.Sprintf("%s:%d", redisHostname, redisPort), "list of clustered redis hosts")

	runBaselineProxyCmd.Flags().BoolVar(&autoRemove, "autoremove", false, "when true, containers are automatically pruned upon exit")
	runBaselineProxyCmd.Flags().StringVar(&logLevel, "log-level", "DEBUG", "log level to set within the running proxy instance")

	runBaselineProxyCmd.Flags().StringVar(&jwtSignerPublicKey, "jwt-signer-public-key", defaultJWTSignerPublicKey, "PEM-encoded public key of the authorized JWT signer for verifying inbound proxy connection attempts")
	runBaselineProxyCmd.Flags().StringVar(&natsAuthToken, "nats-auth-token", "testtoken", "authorization token for the NATS service; will be passed as the -auth argument to NATS")

	runBaselineProxyCmd.Flags().StringVar(&identAPIHost, "ident-host", "ident.provide.services", "hostname of the ident service")
	runBaselineProxyCmd.Flags().StringVar(&identAPIScheme, "ident-scheme", "https", "protocol scheme of the ident service")

	runBaselineProxyCmd.Flags().StringVar(&nchainAPIHost, "nchain-host", "nchain.provide.services", "hostname of the nchain service")
	runBaselineProxyCmd.Flags().StringVar(&nchainAPIScheme, "nchain-scheme", "https", "protocol scheme of the nchain service")

	runBaselineProxyCmd.Flags().StringVar(&privacyAPIHost, "privacy-host", "privacy.provide.services", "hostname of the privacy service")
	runBaselineProxyCmd.Flags().StringVar(&privacyAPIScheme, "privacy-scheme", "https", "protocol scheme of the privacy service")

	runBaselineProxyCmd.Flags().StringVar(&sorID, "sor", "", "primary internal system of record identifier being baselined")
	runBaselineProxyCmd.Flags().StringVar(&sorURL, "sor-url", "https://", "url of the primary internal system of record being baselined")

	// runBaselineProxyCmd.Flags().StringVar(&serviceNowAPIHost, "servicenow-api-host", "", "hostname of the ServiceNow service")
	// runBaselineProxyCmd.Flags().StringVar(&serviceNowAPIScheme, "servicenow-api-scheme", "https", "protocol scheme of the ServiceNow service")
	// runBaselineProxyCmd.Flags().StringVar(&serviceNowAPIPath, "servicenow-api-path", "api/now/table", "base path of the ServiceNow API")
	// runBaselineProxyCmd.Flags().StringVar(&serviceNowAPIUsername, "servicenow-api-username", "", "username to use for basic authorization against the ServiceNow API")
	// runBaselineProxyCmd.Flags().StringVar(&serviceNowAPIPassword, "servicenow-api-password", "", "password to use for basic authorization against the ServiceNow API")

	runBaselineProxyCmd.Flags().StringVar(&vaultAPIHost, "vault-host", "vault.provide.services", "hostname of the vault service")
	runBaselineProxyCmd.Flags().StringVar(&vaultAPIScheme, "vault-scheme", "https", "protocol scheme of the vault service")
	runBaselineProxyCmd.Flags().StringVar(&vaultRefreshToken, "vault-refresh-token", os.Getenv("VAULT_REFRESH_TOKEN"), "refresh token to vend access tokens for use with vault")
	runBaselineProxyCmd.Flags().StringVar(&vaultSealUnsealKey, "vault-seal-unseal-key", os.Getenv("VAULT_SEAL_UNSEAL_KEY"), "seal/unseal key for the vault service")

	runBaselineProxyCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	runBaselineProxyCmd.MarkFlagRequired("organization")

	runBaselineProxyCmd.Flags().StringVar(&organizationRefreshToken, "organization-refresh-token", os.Getenv("PROVIDE_ORGANIZATION_REFRESH_TOKEN"), "refresh token to vend access tokens for use with the local organization")

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

	runBaselineProxyCmd.Flags().StringVar(&baselineOrganizationAddress, "organization-address", defaultBaselineOrganizationAddress, "public baseline regsitry address of the organization")
	runBaselineProxyCmd.Flags().StringVar(&baselineRegistryContractAddress, "registry-contract-address", defaultBaselineRegistryContractAddress, "public baseline regsitry contract address")
	runBaselineProxyCmd.Flags().StringVar(&nchainBaselineNetworkID, "nchain-network-id", defaultNChainBaselineNetworkID, "nchain network id of the baseline mainnet")
}
