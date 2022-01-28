package node

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
	"github.com/provideplatform/provide-cli/cmd/common"

	"github.com/spf13/cobra"
)

const baseledgerContainerImage = "baseledger/node"
const rpcContainerPort = 1337

type portMapping struct {
	hostPort      int
	containerPort int
}

var autoRemove bool
var dockerNetworkID string
var name string

var baseledgerABCIConnectionType string
var baseledgerBlockTime string
var baseledgerBootstrapPeers string
var baseledgerChainID string
var baseledgerFastSync bool
var baseledgerFastSyncVersion string
var baseledgerFilterPeers bool
var baseledgerLogFormat string
var baseledgerLogLevel string
var baseledgerMode string
var baseledgerDBBackend string
var baseledgerGenesisStateURL string
var baseledgerGenesisURL string
var baseledgerMempoolCacheSize int
var baseledgerMempoolSize int
var baseledgerNetworkName string
var baseledgerP2PListenAddress string
var baseledgerP2PMaxConnections int
var baseledgerP2PMaxPacketMessagePayloadSize int
var baseledgerP2PPersistentPeerMaxDialPeriod int
var baseledgerPeerAlias string
var baseledgerPeerBroadcastAddress string
var baseledgerPersistentPeers string
var baseledgerRPCCORSOrigins string
var baseledgerRPCHostname string
var baseledgerRPCListenAddress string
var baseledgerRPCMaxSubscriptionClients int
var baseledgerRPCMaxClientSubscriptions int
var baseledgerRPCMaxOpenConnections int
var baseledgerRPCPort int
var baseledgerSeeds string
var baseledgerStakingContractAddress string
var baseledgerStakingNetwork string
var baseledgerTxIndexer string

var logLevel string
var syslogEndpoint string

var identAPIHost string
var identAPIScheme string

var natsURL string

var nchainAPIHost string
var nchainAPIScheme string

var privacyAPIHost string
var privacyAPIScheme string

var provideRefreshToken string

var vaultAPIHost string
var vaultAPIScheme string

var vaultID string
var vaultKeyID string
var vaultRefreshToken string

var startBaseledgerNodeCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a baseledger node",
	Long:  `Start a local baseledger validator, full or seed node`,
	Run:   startBaseledgerNode,
}

func startBaseledgerNode(cmd *cobra.Command, args []string) {
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to initialize docker; %s", err.Error())
		os.Exit(1)
	}

	go common.PurgeContainers(docker, name, false)

	wg := &sync.WaitGroup{}

	images := make([]string, 0)
	images = append(
		images,
		baseledgerContainerImage,
	)

	for _, image := range images {
		img := image
		wg.Add(1)
		go func() {
			err := pullImage(docker, img)
			if err != nil {
				log.Printf("failed to pull local baseledger container image: %s; %s", img, err.Error())
				os.Exit(1)
			}
			wg.Done()
		}()
	}

	configureNetwork(docker)

	// run baseledger
	wg.Add(1)
	go runBaseledger(docker, wg)

	wg.Wait()
	log.Printf("%s local baseledger instance started", name)
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

func containerEnvironmentFactory(listenPort *int) []string {
	env := make([]string, 0)
	for _, envvar := range []string{
		fmt.Sprintf("BASELEDGER_ABCI_CONNECTION_TYPE=%s", baseledgerABCIConnectionType),
		fmt.Sprintf("BASELEDGER_BLOCK_TIME=%s", baseledgerBlockTime),
		fmt.Sprintf("BASELEDGER_BOOTSTRAP_PEERS=%s", baseledgerBootstrapPeers),
		fmt.Sprintf("BASELEDGER_CHAIN_ID=%s", baseledgerChainID),
		fmt.Sprintf("BASELEDGER_DB_BACKEND=%s", baseledgerDBBackend),
		fmt.Sprintf("BASELEDGER_FAST_SYNC=%v", baseledgerFastSync),
		fmt.Sprintf("BASELEDGER_FAST_SYNC_VERSION=%v", baseledgerFastSyncVersion),
		fmt.Sprintf("BASELEDGER_FILTER_PEERS=%v", baseledgerFilterPeers),
		fmt.Sprintf("BASELEDGER_GENESIS_URL=%s", baseledgerGenesisURL),
		fmt.Sprintf("BASELEDGER_GENESIS_STATE_URL=%s", baseledgerGenesisStateURL),
		fmt.Sprintf("BASELEDGER_LOG_FORMAT=%s", baseledgerLogFormat),
		fmt.Sprintf("BASELEDGER_LOG_LEVEL=%s", baseledgerLogLevel),
		fmt.Sprintf("BASELEDGER_MEMPOOL_SIZE=%d", baseledgerMempoolSize),
		fmt.Sprintf("BASELEDGER_MEMPOOL_CACHE_SIZE=%d", baseledgerMempoolCacheSize),
		fmt.Sprintf("BASELEDGER_MODE=%s", baseledgerMode),
		fmt.Sprintf("BASELEDGER_NETWORK_NAME=%s", baseledgerNetworkName),
		fmt.Sprintf("BASELEDGER_P2P_LISTEN_ADDRESS=%s", baseledgerP2PListenAddress),
		fmt.Sprintf("BASELEDGER_P2P_MAX_CONNECTIONS=%d", baseledgerP2PMaxConnections),
		fmt.Sprintf("BASELEDGER_P2P_MAX_PACKET_MESSAGE_PAYLOAD_SIZE=%d", baseledgerP2PMaxPacketMessagePayloadSize),
		fmt.Sprintf("BASELEDGER_P2P_PERSISTENT_PEER_MAX_DIAL_PERIOD=%d", baseledgerP2PPersistentPeerMaxDialPeriod),
		fmt.Sprintf("BASELEDGER_PEER_ALIAS=%s", baseledgerPeerAlias),
		fmt.Sprintf("BASELEDGER_PEER_BROADCAST_ADDRESS=%s", baseledgerPeerBroadcastAddress),
		fmt.Sprintf("BASELEDGER_PERSISTENT_PEERS=%s", baseledgerPersistentPeers),
		fmt.Sprintf("BASELEDGER_RPC_LISTEN_ADDRESS=%s", baseledgerRPCListenAddress),
		fmt.Sprintf("BASELEDGER_RPC_CORS_ORIGINS=%s", baseledgerRPCCORSOrigins),
		fmt.Sprintf("BASELEDGER_RPC_MAX_OPEN_CONNECTIONS=%d", baseledgerRPCMaxOpenConnections),
		fmt.Sprintf("BASELEDGER_RPC_MAX_SUBSCRIPTION_CLIENTS=%d", baseledgerRPCMaxSubscriptionClients),
		fmt.Sprintf("BASELEDGER_RPC_MAX_CLIENT_SUBSCRIPTIONS=%d", baseledgerRPCMaxClientSubscriptions),
		fmt.Sprintf("BASELEDGER_SEEDS=%s", baseledgerSeeds),
		fmt.Sprintf("BASELEDGER_STAKING_CONTRACT_ADDRESS=%s", baseledgerStakingContractAddress),
		fmt.Sprintf("BASELEDGER_STAKING_NETWORK=%s", baseledgerStakingNetwork),
		fmt.Sprintf("BASELEDGER_TX_INDEXER=%s", baseledgerTxIndexer),

		fmt.Sprintf("LOG_LEVEL=%s", logLevel),
		fmt.Sprintf("NATS_URL=%s", natsURL),
		fmt.Sprintf("PROVIDE_REFRESH_TOKEN=%s", provideRefreshToken),

		fmt.Sprintf("VAULT_ID=%s", vaultID),
		fmt.Sprintf("VAULT_KEY_ID=%s", vaultKeyID),
		fmt.Sprintf("VAULT_REFRESH_TOKEN=%s", vaultRefreshToken),
	} {
		env = append(env, envvar)
	}

	if listenPort != nil {
		env = append(env, fmt.Sprintf("PORT=%d", *listenPort))
	}

	return env
}

func runBaseledger(docker *client.Client, wg *sync.WaitGroup) {
	_, err := runContainer(
		docker,
		fmt.Sprintf("%s-api", strings.ReplaceAll(name, " ", "")),
		baseledgerRPCHostname,
		baseledgerContainerImage,
		&[]string{"./.bin/node"},
		nil,
		&[]string{"CMD", "curl", "-f", fmt.Sprintf("http://%s:%d/status", baseledgerRPCHostname, rpcContainerPort)},
		map[string]string{},
		[]portMapping{{
			hostPort:      baseledgerRPCPort,
			containerPort: rpcContainerPort,
		}}...,
	)

	if err != nil {
		log.Printf("failed to create local baseledger node container; %s", err.Error())
		os.Exit(1)
	}

	if wg != nil {
		wg.Done()
	}
}

func pullImage(docker *client.Client, image string) error {
	log.Printf("pulling local baseledger container image: %s", image)
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
	mounts map[string]string,
	ports ...portMapping,
) (*container.ContainerCreateCreatedBody, error) {
	log.Printf("running local baseledger container image: %s", image)
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

	containerConfig := &container.Config{
		Env:      containerEnvironmentFactory(listenPort),
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

func init() {
	startBaseledgerNodeCmd.Flags().BoolVar(&autoRemove, "autoremove", false, "when true, containers are automatically pruned upon exit")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerABCIConnectionType, "abci-connection-type", "socket", "application blockchain interface connection type for use with tendermint")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerBlockTime, "block-time", "5s", "block time")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerBootstrapPeers, "bootstrap-peers", "", "comma-delimited list of bootstrap peers")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerChainID, "chain-id", "peachtree", "baseledger chain id")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerDBBackend, "db-backend", "peachtree", "baseledger database backend")
	startBaseledgerNodeCmd.Flags().BoolVar(&baseledgerFastSync, "fast-sync", true, "when true, block synchronization and commit verification is parallelized")
	startBaseledgerNodeCmd.Flags().BoolVar(&baseledgerFilterPeers, "filter-peers", false, "when true, baseledger network peers are filtered by way of delegation to the ABCI")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerGenesisStateURL, "genesis-state-url", "", "url of the genesis state JSON")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerGenesisURL, "genesis-url", "http://genesis.peachtree.baseledger.provide.network:1337/genesis", "url of the network genesis JSON; fetched via peer RPC if left blank")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerLogFormat, "tendermint-log-format", "plain", "log format to set within the local baseledger node")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerLogLevel, "tendermint-log-level", "main:debug,abci-client:debug,blockchain:debug,consensus:debug,state:debug,statesync:debug,*:error", "log level to set within the local baseledger node")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerMempoolCacheSize, "mempool-cache-size", 256, "number of cached transactions to allow in the mempool at any given time")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerMempoolSize, "mempool-size", 1024, "number of transactions to allow in the mempool at any given time")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerMode, "mode", "full", "mode in which to run the baseledger node (i.e., validator, full or seed)")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerP2PListenAddress, "p2p-listen-address", "tcp://0.0.0.0:33333", "peer-to-peer listen address")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerP2PMaxConnections, "p2p-max-connections", 32, "maximum number of inbound and outbound peer-to-peer connections")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerP2PMaxPacketMessagePayloadSize, "p2p-max-message-packet-payload-size", 22020096, "maximum size, in bytes, of peer-to-peer message packets")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerP2PPersistentPeerMaxDialPeriod, "p2p-max-persistent-peer-dial-period", int(time.Second*10), "maximum pause when redialing a persistent peer (if zero, exponential backoff is used)")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerPeerAlias, "peer-alias", "prvd", "node alias to advertise to other peers in the network")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerPeerBroadcastAddress, "peer-broadcast-address", "", "address to advertise to other nodes in the network")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerPersistentPeers, "persistent-peers", "", "comma-delimited list of persistent peers")
	startBaseledgerNodeCmd.Flags().StringVar(&natsURL, "nats-url", "nats://35.174.77.159:4222,nats://35.172.123.165:4222", "NATS cluster url on which to receive staking contract events")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerRPCCORSOrigins, "rpc-cors-origins", "*", "CORS origins for the local baseledger RPC service")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerRPCHostname, "rpc-hostname", fmt.Sprintf("%s-api", name), "hostname for the local baseledger RPC service")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerRPCListenAddress, "rpc-listen-address", "tcp://0.0.0.0:1337", "listen address for the local baseledger RPC service")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerRPCMaxOpenConnections, "rpc-max-open-connections", 1024, "maximum number of open RPC connections")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerRPCPort, "rpc-port", 1337, "host port on which to expose the local baseledger RPC service")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerRPCMaxSubscriptionClients, "rpc-max-subscription-clients", 1024, "maximum number of concurrent subscription clients")
	startBaseledgerNodeCmd.Flags().IntVar(&baseledgerRPCMaxClientSubscriptions, "rpc-max-client-subscriptions", 32, "maximum number of subscriptions allowed per client")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerSeeds, "seeds", "", "comma-delimited list of seed nodes")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerStakingContractAddress, "staking-contract-address", "", "address of the staking contract on the named --staking-network")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerStakingNetwork, "staking-network", "ropsten", "name of the staking network")
	startBaseledgerNodeCmd.Flags().StringVar(&syslogEndpoint, "syslog-endpoint", "", "syslog endpoint to which syslog udp packets will be sent")
	startBaseledgerNodeCmd.Flags().StringVar(&baseledgerTxIndexer, "tx-indexer", "kv", "transaction indexing engine")

	startBaseledgerNodeCmd.Flags().StringVar(&identAPIHost, "ident-host", "ident.provide.services", "hostname of the ident service")
	startBaseledgerNodeCmd.Flags().StringVar(&identAPIScheme, "ident-scheme", "https", "protocol scheme of the ident service")

	startBaseledgerNodeCmd.Flags().StringVar(&logLevel, "log-level", "DEBUG", "log level to use in the baseledger node outside of tendermint")
	startBaseledgerNodeCmd.Flags().StringVar(&name, "name", "baseledger-local", "name of the baseledger node instance")

	startBaseledgerNodeCmd.Flags().StringVar(&nchainAPIHost, "nchain-host", "nchain.provide.services", "hostname of the nchain service")
	startBaseledgerNodeCmd.Flags().StringVar(&nchainAPIScheme, "nchain-scheme", "https", "protocol scheme of the nchain service")

	startBaseledgerNodeCmd.Flags().StringVar(&privacyAPIHost, "privacy-host", "privacy.provide.services", "hostname of the privacy service")
	startBaseledgerNodeCmd.Flags().StringVar(&privacyAPIScheme, "privacy-scheme", "https", "protocol scheme of the privacy service")

	startBaseledgerNodeCmd.Flags().StringVar(&provideRefreshToken, "provide-refresh-token", "", "refresh token to vend access tokens for use with the Provide stack")

	startBaseledgerNodeCmd.Flags().StringVar(&vaultAPIHost, "vault-host", "vault.provide.services", "hostname of the vault service")
	startBaseledgerNodeCmd.Flags().StringVar(&vaultAPIScheme, "vault-scheme", "https", "protocol scheme of the vault service")
	startBaseledgerNodeCmd.Flags().StringVar(&vaultRefreshToken, "vault-refresh-token", os.Getenv("VAULT_REFRESH_TOKEN"), "refresh token to vend access tokens for use with vault")
}
