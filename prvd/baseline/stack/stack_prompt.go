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
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/spf13/cobra"
)

const promptStepStart = "Start"
const promptStepStop = "Stop"
const promptStepLogs = "Logs"

var emptyPromptArgs = []string{promptStepStart, promptStepStop, promptStepLogs}
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

var SoRPromptArgs = []string{"SAP", "ServiceNow", "Sales Force"}
var SoRPromptLabel = "Select a Sor"

// General Endpoints
func generalPrompt(cmd *cobra.Command, args []string, currentStep string) {
	switch step := currentStep; step {
	case promptStepStart:
		common.RequireOrganization()
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
				sorURL = common.FreeInput("System of Record URL", "", common.NoValidation)
			}
			if apiHostname == "" {
				apiHostname = common.FreeInput("API Hostname", "", common.NoValidation)
			}
			if port == 8080 {
				port, _ = strconv.Atoi(common.FreeInput("Port", "8080", common.NumberValidation))
			}
			if consumerHostname == name+"-consumer" {
				consumerHostname = common.FreeInput("Consumer Hostname", name+"-consumer", common.NoValidation)
			}
			if natsHostname == name+"-nats" {
				natsHostname = common.FreeInput("NATS Hostname", name+"-nats", common.NoValidation)
			}
			if natsPort == 4222 {
				natsPort, _ = strconv.Atoi(common.FreeInput("NATS Port", "4222", common.NumberValidation))
			}
			if natsWebsocketPort == 4221 {
				natsWebsocketPort, _ = strconv.Atoi(common.FreeInput("NATS Websocket Port", "4221", common.NumberValidation))
			}
			if natsAuthToken == "testtoken" {
				natsAuthToken = common.FreeInput("NATS Auth Token", "testtoken", common.NoValidation)
			}
			if redisHostname == fmt.Sprintf("%s-redis", name) {
				redisHostname = common.FreeInput("Redis Host Name", name+"-redis", common.NoValidation)
			}
			if redisPort == 6379 {
				redisPort, _ = strconv.Atoi(common.FreeInput("Redis Port", "6379", common.NumberValidation))
			}
			if redisHosts == redisHostname+":"+strconv.Itoa(redisContainerPort) {
				redisPort, _ = strconv.Atoi(common.FreeInput("Redis Port", redisHostname+":"+strconv.Itoa(redisContainerPort), common.NoValidation))
			}
			if !autoRemove {
				autoRemove = common.SelectInput(boolPromptArgs, autoRemovePromptLabel) == "Yes"
			}
			if strings.ToLower(logLevel) == "debug" {
				logLevel = common.FreeInput("Log Level", "debug", common.NoValidation)
			}
			if jwtSignerPublicKey == "" {
				jwtSignerPublicKey = common.FreeInput("JWT Signer Public Key", "", common.NoValidation)
			}
			if identAPIHost == "ident.provide.services" {
				nchainAPIHost = common.FreeInput("Ident API Host", "ident.provide.services", common.NoValidation)
			}
			if identAPIScheme == "https" {
				nchainAPIScheme = common.FreeInput("Ident API Scheme", "https", common.NoValidation)
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
				baselineOrganizationAddress = common.FreeInput("Baseline Organization Address", "0x", common.NoValidation)
			}
			if baselineRegistryContractAddress == "0x" {
				baselineOrganizationAddress = common.FreeInput("Baseline Registry Contract Address", "0x", common.HexValidation)
			}
			if common.WorkgroupID == "" {
				baselineOrganizationAddress = common.FreeInput("Baseline Workgroup ID", "", common.HexValidation)
			}
			if nchainBaselineNetworkID == "0x" {
				baselineOrganizationAddress = common.FreeInput("Nchain Baseline Network ID", "0x", common.HexValidation)
			}
		}
		runStackStart(cmd, args)
	case promptStepStop:
		if Optional {
			fmt.Println("Optional Flags:")
			if name == "" {
				name = common.FreeInput("Name", "", common.NoValidation)
			}
		}
		runStackStop(cmd, args)
	case promptStepLogs:
		if Optional {
			fmt.Println("Optional Flags:")
			if name == "" {
				name = common.FreeInput("Name", "", common.NoValidation)
			}
		}
		stackLogsRun(cmd, args)
	case "":
		result := common.SelectInput(emptyPromptArgs, emptyPromptLabel)
		generalPrompt(cmd, args, result)
	}
}
