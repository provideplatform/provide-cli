package stack

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/spf13/cobra"
)

const promptStepRun = "Run"
const promptStepStop = "Stop"
const promptStepLogs = "Logs"

var emptyPromptArgs = []string{promptStepRun, promptStepStop, promptStepLogs}
var emptyPromptLabel = "What would you like to do"

var boolPromptArgs = []string{"No", "Yes"}
var tunnelPromptLabel = "Would you like to set up a tunnel"
var tunnelAPIPromptLabel = "Would you like to set up a API tunnel"
var tunnelMessagingPromptLabel = "Would you like to set up a messaging tunnel"
var autoRemovePromptLabel = "Would you like to automatically remove"
var localVaultPromptLabel = "Would you like to set up vault locally"
var localIdentPromptLabel = "Would you like to set up ident locally"
var localNchainPromptLabel = "Would you like to set up nachain locally"
var localPrivacyPromptLabel = "Would you like to set up privacy locally"

var SoRPromptArgs = []string{"SAP", "Service Now", "Sales Force"}
var SoRPromptLabel = "Select a Sor"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepRun:
		common.RequireOrganization()
		if common.ManifestLoad() {
			setVarsContents()
		}
		if Optional {
			if name == "" {
				name = common.FreeInput("Name", "", common.NoValidation)
			}
			if common.APIEndpoint == "" {
				common.APIEndpoint = common.FreeInput("API endpoint", "", common.NoValidation)
			}
			if common.MessagingEndpoint == "" {
				common.MessagingEndpoint = common.FreeInput("Messaging endpoint", "", common.NoValidation)
			}
			if !common.Tunnel {
				common.Tunnel = common.SelectInput(boolPromptArgs, tunnelPromptLabel) == "Yes"
			}
			if !common.ExposeAPITunnel {
				common.ExposeAPITunnel = common.SelectInput(boolPromptArgs, tunnelAPIPromptLabel) == "Yes"
			}
			if !common.ExposeMessagingTunnel {
				common.ExposeMessagingTunnel = common.SelectInput(boolPromptArgs, tunnelMessagingPromptLabel) == "Yes"
			}
			if sorID == "" {
				sorID = common.SelectInput(SoRPromptArgs, SoRPromptLabel)
			}
			if sorURL == "" {
				sorURL = common.FreeInput("SoR URL", "", common.NoValidation)
			}
			if apiHostname == "" {
				apiHostname = common.FreeInput("API Hostname", "", common.NoValidation)
			}
			if port == apiContainerPort {
				port, _ = strconv.Atoi(common.FreeInput("Port", "8080", common.NumberValidation))
			}
			if consumerHostname == name+"-consumer" {
				consumerHostname = common.FreeInput("Consumer Hostname", name+"-consumer", common.NoValidation)
			}
			if natsHostname == name+"-nats" {
				natsHostname = common.FreeInput("Nats Hostname", name+"-nats", common.NoValidation)
			}
			if natsPort == natsContainerPort {
				natsPort, _ = strconv.Atoi(common.FreeInput("Nats Port", "4222", common.NumberValidation))
			}
			if natsWebsocketPort == natsWebsocketContainerPort {
				natsWebsocketPort, _ = strconv.Atoi(common.FreeInput("Nats Websocket Port", "4221", common.NumberValidation))
			}
			if natsAuthToken == "testtoken" {
				natsAuthToken = common.FreeInput("Nats Auth Token", "testtoken", common.NoValidation)
			}
			if natsStreamingHostname == name+"-nats-streaming" {
				natsStreamingHostname = common.FreeInput("Nats Streaming Token", name+"-nats-streaming", common.NoValidation)
			}
			if natsStreamingPort == natsStreamingContainerPort {
				natsStreamingPort, _ = strconv.Atoi(common.FreeInput("Nats Streaming Port", "4221", common.NumberValidation))
			}
			if redisHostname == name+"-reddis" {
				redisHostname = common.FreeInput("Reddis Host Name", name+"-reddis", common.NoValidation)
			}
			if redisPort == redisContainerPort {
				redisPort, _ = strconv.Atoi(common.FreeInput("Reddis Port", "6379", common.NumberValidation))
			}
			if redisHosts == redisHostname+":"+strconv.Itoa(redisContainerPort) {
				redisHosts = common.FreeInput("Reddis Port", redisHostname+":"+strconv.Itoa(redisContainerPort), common.NoValidation)
			}
			if !autoRemove {
				autoRemove = common.SelectInput(boolPromptArgs, autoRemovePromptLabel) == "Yes"
			}
			if logLevel == "DEBUG" {
				logLevel = common.FreeInput("Reddis Host Name", "DEBUG", common.NoValidation)
			}
			if jwtSignerPublicKey == "" {
				jwtSignerPublicKey = common.FreeInput("JWT Signer Public Key", "", common.NoValidation)
			}
			if identAPIHost == "ident.provide.services" {
				identAPIHost = common.FreeInput("Ident API Host", "ident.provide.services", common.NoValidation)
			}
			if identAPIScheme == "https" {
				identAPIScheme = common.FreeInput("Ident API Scheme", "https", common.NoValidation)
			}
			if nchainAPIHost == "nchain.provide.services" {
				nchainAPIHost = common.FreeInput("Nchain API Host", "nchain.provide.services", common.NoValidation)
			}
			if nchainAPIScheme == "https" {
				nchainAPIScheme = common.FreeInput("Nchain API Scheme", "https", common.NoValidation)
			}
			if privacyAPIHost == "privacy.provide.services" {
				privacyAPIHost = common.FreeInput("Privacy API Host", "privacy.provide.services", common.NoValidation)
			}
			if privacyAPIScheme == "https" {
				privacyAPIScheme = common.FreeInput("Privacy API Scheme", "https", common.NoValidation)
			}
			if vaultAPIHost == "vault.provide.services" {
				vaultAPIHost = common.FreeInput("Vault API Host", "vault.provide.services", common.NoValidation)
			}
			if vaultAPIScheme == "https" {
				vaultAPIScheme = common.FreeInput("Vault API Scheme", "https", common.NoValidation)
			}
			if vaultRefreshToken == os.Getenv("VAULT_REFRESH_TOKEN") {
				vaultRefreshToken = common.FreeInput("Vault API Refresh Token", os.Getenv("VAULT_REFRESH_TOKEN"), common.NoValidation)
			}
			if vaultSealUnsealKey == os.Getenv("VAULT_SEAL_UNSEAL_KEY") {
				vaultSealUnsealKey = common.FreeInput("Vault Un/Seal Token", os.Getenv("VAULT_SEAL_UNSEAL_KEY"), common.NoValidation)
			}
			if !withLocalVault {
				withLocalVault = strings.ToLower(common.SelectInput(boolPromptArgs, localVaultPromptLabel)) == "yes"
			}
			if !withLocalIdent {
				withLocalIdent = strings.ToLower(common.SelectInput(boolPromptArgs, localIdentPromptLabel)) == "yes"
			}
			if !withLocalNChain {
				withLocalNChain = strings.ToLower(common.SelectInput(boolPromptArgs, localNchainPromptLabel)) == "yes"
			}
			if !withLocalPrivacy {
				withLocalPrivacy = strings.ToLower(common.SelectInput(boolPromptArgs, localPrivacyPromptLabel)) == "yes"
			}
			if organizationRefreshToken == os.Getenv("PROVIDE_ORGANIZATION_REFRESH_TOKEN") {
				organizationRefreshToken = common.FreeInput("Organization Refresh Token", os.Getenv("PROVIDE_ORGANIZATION_REFRESH_TOKEN"), common.NoValidation)
			}
			if baselineOrganizationAddress == "0x" {
				baselineOrganizationAddress = common.FreeInput("Baseline Organization Address", "0x", common.HexValidation)
			}
			if baselineRegistryContractAddress == "0x" {
				baselineRegistryContractAddress = common.FreeInput("Baseline Registry Contract Address", "0x", common.HexValidation)
			}
			if baselineWorkgroupID == "" {
				baselineWorkgroupID = common.FreeInput("Baseline Workgroup ID", "", common.HexValidation)
			}
			if nchainBaselineNetworkID == "0x" {
				nchainBaselineNetworkID = common.FreeInput("Nchain Baseline Network ID", "0x", common.HexValidation)
			}
			common.ManifestSave(marshalEnvManifest())

		}
		runProxyRun(cmd, args)
	case promptStepStop:
		if Optional {
			fmt.Println("Optional Flags:")
			if name == "" {
				name = common.FreeInput("Name", "", common.NoValidation)
			}
		}
		stopProxyRun(cmd, args)
	case promptStepLogs:
		if Optional {
			fmt.Println("Optional Flags:")
			if name == "" {
				name = common.FreeInput("Name", "", common.NoValidation)
			}
		}
		logsProxyRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}

// var portMappings = common.PortMapping{
// 	[]portMapping{
// 		hostPort:      port,
// 		containerPort: apiContainerPort,
// 	},
// }
func marshalEnvManifest() []byte {
	content, _ := json.MarshalIndent(common.EnvManifest{
		APIEndpoint:                      common.APIEndpoint,
		MessagingEndpoint:                common.MessagingEndpoint,
		Tunnel:                           common.Tunnel,
		ExposeAPITunnel:                  common.ExposeAPITunnel,
		ExposeMessagingTunnel:            common.ExposeMessagingTunnel,
		ApiContainerPort:                 apiContainerPort,
		NatsContainerPort:                natsContainerPort,
		NatsWebsocketContainerPort:       natsWebsocketContainerPort,
		NatsStreamingContainerPort:       natsStreamingContainerPort,
		PostgresContainerPort:            postgresContainerPort,
		RedisContainerPort:               redisContainerPort,
		DockerNetworkID:                  dockerNetworkID,
		Name:                             name,
		Port:                             port,
		IdentPort:                        identPort,
		NchainPort:                       nchainPort,
		PrivacyPort:                      privacyPort,
		VaultPort:                        vaultPort,
		NatsPort:                         natsPort,
		NatsWebsocketPort:                natsWebsocketPort,
		NatsStreamingPort:                natsStreamingPort,
		PostgresPort:                     postgresPort,
		RedisPort:                        redisPort,
		ApiHostname:                      apiHostname,
		ConsumerHostname:                 consumerHostname,
		IdentAPIHost:                     identAPIHost,
		IdentHostname:                    identHostname,
		IdentConsumerHostname:            identConsumerHostname,
		NchainHostname:                   nchainHostname,
		NchainConsumerHostname:           nchainConsumerHostname,
		NchainStatsdaemonHostname:        nchainStatsdaemonHostname,
		NchainReachabilitydaemonHostname: nchainReachabilitydaemonHostname,
		PrivacyHostname:                  privacyHostname,
		PrivacyConsumerHostname:          privacyConsumerHostname,
		VaultHostname:                    vaultHostname,
		NatsHostname:                     natsHostname,
		NatsServerName:                   natsServerName,
		NatsStreamingHostname:            natsStreamingHostname,
		PostgresHostname:                 postgresHostname,
		RedisHostname:                    redisHostname,
		RedisHosts:                       redisHosts,
		AutoRemove:                       autoRemove,
		LogLevel:                         logLevel,
		BaselineOrganizationAddress:      baselineOrganizationAddress,
		BaselineRegistryContractAddress:  baselineRegistryContractAddress,
		BaselineWorkgroupID:              baselineWorkgroupID,
		NchainBaselineNetworkID:          nchainBaselineNetworkID,
		JwtSignerPublicKey:               jwtSignerPublicKey,
		NatsAuthToken:                    natsAuthToken,
		IdentAPIScheme:                   identAPIScheme,
		NchainAPIHost:                    nchainAPIHost,
		NchainAPIScheme:                  nchainAPIScheme,
		WorkgroupAccessToken:             workgroupAccessToken,
		OrganizationRefreshToken:         organizationRefreshToken,
		PrivacyAPIHost:                   privacyAPIHost,
		PrivacyAPIScheme:                 privacyAPIScheme,
		SorID:                            sorID,
		SorURL:                           sorURL,
		VaultAPIHost:                     vaultAPIHost,
		VaultAPIScheme:                   vaultAPIScheme,
		VaultRefreshToken:                vaultRefreshToken,
		VaultSealUnsealKey:               vaultSealUnsealKey,
		SapAPIHost:                       sapAPIHost,
		SapAPIScheme:                     sapAPIScheme,
		SapAPIUsername:                   sapAPIUsername,
		SapAPIPassword:                   sapAPIPassword,
		SapAPIPath:                       sapAPIPath,
		ServiceNowAPIHost:                serviceNowAPIHost,
		ServiceNowAPIScheme:              serviceNowAPIScheme,
		ServiceNowAPIUsername:            serviceNowAPIUsername,
		ServiceNowAPIPassword:            serviceNowAPIPassword,
		ServiceNowAPIPath:                serviceNowAPIPath,
		SalesforceAPIHost:                salesforceAPIHost,
		SalesforceAPIScheme:              salesforceAPIScheme,
		SalesforceAPIPath:                salesforceAPIPath,
		WithLocalVault:                   withLocalVault,
		WithLocalIdent:                   withLocalIdent,
		WithLocalNChain:                  withLocalNChain,
		WithLocalPrivacy:                 withLocalPrivacy,
	}, "", "  ")
	return content
}

func setVarsContents() {
	common.APIEndpoint = common.LoadedManifest.APIEndpoint
	common.MessagingEndpoint = common.LoadedManifest.MessagingEndpoint
	common.Tunnel = common.LoadedManifest.Tunnel
	common.ExposeAPITunnel = common.LoadedManifest.ExposeAPITunnel
	common.ExposeMessagingTunnel = common.LoadedManifest.ExposeMessagingTunnel
	dockerNetworkID = common.LoadedManifest.DockerNetworkID
	name = common.LoadedManifest.Name
	port = common.LoadedManifest.Port
	identPort = common.LoadedManifest.IdentPort
	nchainPort = common.LoadedManifest.NchainPort
	privacyPort = common.LoadedManifest.PrivacyPort
	vaultPort = common.LoadedManifest.VaultPort
	natsPort = common.LoadedManifest.NatsPort
	natsWebsocketPort = common.LoadedManifest.NatsWebsocketPort
	natsStreamingPort = common.LoadedManifest.NatsStreamingPort
	postgresPort = common.LoadedManifest.PostgresPort
	redisPort = common.LoadedManifest.RedisPort
	apiHostname = common.LoadedManifest.ApiHostname
	consumerHostname = common.LoadedManifest.ConsumerHostname
	identAPIHost = common.LoadedManifest.IdentAPIHost
	identHostname = common.LoadedManifest.IdentHostname
	identConsumerHostname = common.LoadedManifest.IdentConsumerHostname
	nchainHostname = common.LoadedManifest.NchainHostname
	nchainConsumerHostname = common.LoadedManifest.NchainConsumerHostname
	nchainStatsdaemonHostname = common.LoadedManifest.NchainStatsdaemonHostname
	nchainReachabilitydaemonHostname = common.LoadedManifest.NchainReachabilitydaemonHostname
	privacyHostname = common.LoadedManifest.PrivacyHostname
	privacyConsumerHostname = common.LoadedManifest.PrivacyConsumerHostname
	vaultHostname = common.LoadedManifest.VaultHostname
	natsHostname = common.LoadedManifest.NatsHostname
	natsServerName = common.LoadedManifest.NatsServerName
	natsStreamingHostname = common.LoadedManifest.NatsStreamingHostname
	postgresHostname = common.LoadedManifest.PostgresHostname
	redisHostname = common.LoadedManifest.RedisHostname
	redisHosts = common.LoadedManifest.RedisHosts
	autoRemove = common.LoadedManifest.AutoRemove
	logLevel = common.LoadedManifest.LogLevel
	baselineOrganizationAddress = common.LoadedManifest.BaselineOrganizationAddress
	baselineRegistryContractAddress = common.LoadedManifest.BaselineRegistryContractAddress
	baselineWorkgroupID = common.LoadedManifest.BaselineWorkgroupID
	nchainBaselineNetworkID = common.LoadedManifest.NchainBaselineNetworkID
	jwtSignerPublicKey = common.LoadedManifest.JwtSignerPublicKey
	natsAuthToken = common.LoadedManifest.NatsAuthToken
	identAPIScheme = common.LoadedManifest.IdentAPIScheme
	nchainAPIHost = common.LoadedManifest.NchainAPIHost
	nchainAPIScheme = common.LoadedManifest.NchainAPIScheme
	workgroupAccessToken = common.LoadedManifest.WorkgroupAccessToken
	organizationRefreshToken = common.LoadedManifest.OrganizationRefreshToken
	privacyAPIHost = common.LoadedManifest.PrivacyAPIHost
	privacyAPIScheme = common.LoadedManifest.PrivacyAPIScheme
	sorID = common.LoadedManifest.SorID
	sorURL = common.LoadedManifest.SorURL
	vaultAPIHost = common.LoadedManifest.VaultAPIHost
	vaultAPIScheme = common.LoadedManifest.VaultAPIScheme
	vaultRefreshToken = common.LoadedManifest.VaultRefreshToken
	vaultSealUnsealKey = common.LoadedManifest.VaultSealUnsealKey
	sapAPIHost = common.LoadedManifest.SapAPIHost
	sapAPIScheme = common.LoadedManifest.SapAPIScheme
	sapAPIUsername = common.LoadedManifest.SapAPIUsername
	sapAPIPassword = common.LoadedManifest.SapAPIPassword
	sapAPIPath = common.LoadedManifest.SapAPIPath
	serviceNowAPIHost = common.LoadedManifest.ServiceNowAPIHost
	serviceNowAPIScheme = common.LoadedManifest.ServiceNowAPIScheme
	serviceNowAPIUsername = common.LoadedManifest.ServiceNowAPIUsername
	serviceNowAPIPassword = common.LoadedManifest.ServiceNowAPIPassword
	serviceNowAPIPath = common.LoadedManifest.ServiceNowAPIPath
	salesforceAPIHost = common.LoadedManifest.SalesforceAPIHost
	salesforceAPIScheme = common.LoadedManifest.SalesforceAPIScheme
	salesforceAPIPath = common.LoadedManifest.SalesforceAPIPath
	withLocalVault = common.LoadedManifest.WithLocalVault
	withLocalIdent = common.LoadedManifest.WithLocalIdent
	withLocalNChain = common.LoadedManifest.WithLocalNChain
	withLocalPrivacy = common.LoadedManifest.WithLocalPrivacy
}
