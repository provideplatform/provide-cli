package stack

import (
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

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepRun:
		common.RequireOrganization()
		if Optional {
			if name == "" {
				name = common.FreeInput("Name", "", "")
			}
			if common.APIEndpoint == "" {
				common.APIEndpoint = common.FreeInput("API endpoint", "", "")
			}
			if common.MessagingEndpoint == "" {
				common.MessagingEndpoint = common.FreeInput("Messaging endpoint", "", "")
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
			// TODO
			if sorID == "" {
				sorID = common.FreeInput("SoR", "", "")
			}
			if sorURL == "" {
				sorURL = common.FreeInput("SoR URL", "", "")
			}
			if apiHostname == "" {
				apiHostname = common.FreeInput("API Hostname", "", "")
			}
			if port == 8080 {
				port, _ = strconv.Atoi(common.FreeInput("Port", "8080", ""))
			}
			if consumerHostname == name+"-consumer" {
				consumerHostname = common.FreeInput("Consumer Hostname", name+"-consumer", "")
			}
			if natsHostname == name+"-nats" {
				natsHostname = common.FreeInput("Nats Hostname", name+"-nats", "")
			}
			if natsPort == 4222 {
				natsPort, _ = strconv.Atoi(common.FreeInput("Nats Port", "4222", ""))
			}
			if natsWebsocketPort == 4221 {
				natsWebsocketPort, _ = strconv.Atoi(common.FreeInput("Nats Websocket Port", "4221", ""))
			}
			if natsAuthToken == "testtoken" {
				natsAuthToken = common.FreeInput("Nats Auth Token", "testtoken", "")
			}
			if natsStreamingHostname == name+"-nats-streaming" {
				natsStreamingHostname = common.FreeInput("Nats Streaming Token", name+"-nats-streaming", "")
			}
			if natsStreamingPort == 4220 {
				natsStreamingPort, _ = strconv.Atoi(common.FreeInput("Nats Streaming Port", "4221", ""))
			}
			if redisHostname == name+"-reddis" {
				redisHostname = common.FreeInput("Reddis Host Name", name+"-reddis", "")
			}
			if redisPort == 6379 {
				redisPort, _ = strconv.Atoi(common.FreeInput("Reddis Port", "6379", ""))
			}
			// if redisHosts == redisHostname+":"+strconv.Atoi(redisContainerPort) {
			// 	redisPort, _ = strconv.Atoi(common.FreeInput("Reddis Port"))
			// }
			if !autoRemove {
				autoRemove = common.SelectInput(boolPromptArgs, autoRemovePromptLabel) == "Yes"
			}
			if logLevel == "DEBUG" {
				logLevel = common.FreeInput("Reddis Host Name", "DEBUG", "")
			}
			if jwtSignerPublicKey == "" {
				jwtSignerPublicKey = common.FreeInput("JWT Signer Public Key", "", "")
			}
			if identAPIHost == "ident.provide.services" {
				nchainAPIHost = common.FreeInput("Ident API Host", "ident.provide.services", "")
			}
			if identAPIScheme == "https" {
				nchainAPIScheme = common.FreeInput("Ident API Scheme", "https", "")
			}
			if nchainAPIHost == "nchain.provide.services" {
				nchainAPIHost = common.FreeInput("Nchain API Host", "nchain.provide.services", "")
			}
			if nchainAPIScheme == "https" {
				nchainAPIScheme = common.FreeInput("Nchain API Scheme", "https", "")
			}
			if privacyAPIHost == "privacy.provide.services" {
				privacyAPIHost = common.FreeInput("Privacy API Host", "privacy.provide.services", "")
			}
			if privacyAPIScheme == "https" {
				privacyAPIScheme = common.FreeInput("Privacy API Scheme", "https", "")
			}
			if vaultAPIHost == "vault.provide.services" {
				vaultAPIHost = common.FreeInput("Vault API Host", "vault.provide.services", "")
			}
			if vaultAPIScheme == "https" {
				vaultAPIScheme = common.FreeInput("Vault API Scheme", "https", "")
			}
			if vaultRefreshToken == os.Getenv("VAULT_REFRESH_TOKEN") {
				vaultRefreshToken = common.FreeInput("Vault API Refresh Token", "VAULT_REFRESH_TOKEN", "")
			}
			if vaultSealUnsealKey == os.Getenv("VAULT_SEAL_UNSEAL_KEY") {
				vaultSealUnsealKey = common.FreeInput("Vault Un/Seal Token", "VAULT_SEAL_UNSEAL_KEY", "")
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
				organizationRefreshToken = common.FreeInput("Organization Refresh Token", "PROVIDE_ORGANIZATION_REFRESH_TOKEN", "")
			}
			if baselineOrganizationAddress == "0x" {
				baselineOrganizationAddress = common.FreeInput("Baseline Organization Address", "0x", "")
			}
			if baselineRegistryContractAddress == "0x" {
				baselineOrganizationAddress = common.FreeInput("Baseline Registry Contract Address", "0x", "")
			}
			if baselineWorkgroupID == "" {
				baselineOrganizationAddress = common.FreeInput("Baseline Workgroup ID", "", "")
			}
			if nchainBaselineNetworkID == "0x" {
				baselineOrganizationAddress = common.FreeInput("Nchain Baseline Network ID", "0x", "")
			}
		}
		runProxyRun(cmd, args)
	case promptStepStop:
		if Optional {
			fmt.Println("Optional Flags:")
			if name == "" {
				name = common.FreeInput("Name", "", "")
			}
		}
		stopProxyRun(cmd, args)
	case promptStepLogs:
		if Optional {
			fmt.Println("Optional Flags:")
			if name == "" {
				name = common.FreeInput("Name", "", "")
			}
		}
		logsProxyRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
