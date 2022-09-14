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

package systems

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/spf13/cobra"
)

var systemType string

var systemName string
var systemDescription string

var systemMiddlewareType string

var systemEndpointURL string

// could combine the inbound and outbound system vars into systemEndpointURL slices similar to systemAuthMethods, systemUsernames, etc
var systemInboundMiddleware string
var systemInboundEndpointURL string

var systemOutboundMiddleware string
var systemOutboundEndpointURL string

var noMiddlewareAuthMethod string
var noMiddlewareUsername string
var noMiddlewarePassword string
var noMiddlewareRequireClientCredentials bool
var noMiddlewareClientID string
var noMiddlewareClientSecret string

var inboundAuthMethod string
var inboundUsername string
var inboundPassword string
var inboundRequireClientCredentials bool
var inboundClientID string
var inboundClientSecret string

var outboundAuthMethod string
var outboundUsername string
var outboundPassword string
var outboundRequireClientCredentials bool
var outboundClientID string
var outboundClientSecret string

const sapSystemIdentifier = "sap"
const servicenowSystemIdentifier = "servicenow"

const systemNoMiddlewareIdentifier = "No Middleware"
const systemInboundOnlyMiddlewareIdentifier = "Inbound Only"
const systemOutboundOnlyMiddlewareIdentifier = "Outbound Only"
const systemInboundAndOutboundMiddlewareIdentifier = "Inbound & Outbound"

const basicAuthMethodIdentifier = "Basic Auth"

var initBaselineSystemCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize baseline system",
	Long:  `Initialize and configure a new baseline system of record`,
	Run:   initSystem,
}

func initSystem(cmd *cobra.Command, args []string) {
	generalPrompt(cmd, args, promptStepInit)
}

func initSystemRun(cmd *cobra.Command, args []string) {
	if common.OrganizationID == "" {
		common.RequireOrganization()
	}
	if common.WorkgroupID == "" {
		common.RequireWorkgroup()
	}
	if systemType == "" {
		systemTypePrompt()
	}

	common.AuthorizeOrganizationContext(true)

	token, err := common.ResolveOrganizationToken()
	if err != nil {
		fmt.Printf("failed to initialize system; %s", err.Error())
		os.Exit(1)
	}

	var params map[string]interface{}
	systemPrompt(&params, *token.AccessToken)

	vaults, err := vault.ListVaults(*token.AccessToken, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to initialize system; %s", err.Error())
		os.Exit(1)
	}

	if len(vaults) == 0 {
		fmt.Printf("failed to initialize system; workgroup must have a vault")
		os.Exit(1)
	}

	subjectAccountID := common.SHA256(fmt.Sprintf("%s.%s", common.OrganizationID, common.WorkgroupID))
	sa, err := baseline.GetSubjectAccountDetails(*token.AccessToken, common.OrganizationID, subjectAccountID, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to initialize system; %s", err.Error())
		os.Exit(1)
	}

	isOnboarded := sa.ID != nil
	if !isOnboarded {
		raw, _ := json.Marshal(params)

		secretParams := map[string]interface{}{
			"type":  systemType,
			"name":  systemName,
			"value": string(raw),
		}

		secret, err := vault.CreateSecret(*token.AccessToken, vaults[0].ID.String(), secretParams)
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		isOperator := common.Organization.Metadata.Workgroups[common.Workgroup.ID].OperatorSeparationDegree == 0
		if isOperator {
			if common.Workgroup.Config.SystemSecretIDs == nil {
				common.Workgroup.Config.SystemSecretIDs = make([]*uuid.UUID, 0)
			}

			common.Workgroup.Config.SystemSecretIDs = append(common.Workgroup.Config.SystemSecretIDs, &secret.ID)

			var wgInterface map[string]interface{}
			raw, _ := json.Marshal(common.Workgroup)
			json.Unmarshal(raw, &wgInterface)

			if err := baseline.UpdateWorkgroup(*token.AccessToken, common.Workgroup.ID.String(), wgInterface); err != nil {
				fmt.Printf("failed to initialize system; %s", err.Error())
				os.Exit(1)
			}
		}

		common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs = append(common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs, &secret.ID)

		var orgInterface map[string]interface{}
		raw, _ = json.Marshal(common.Organization)
		json.Unmarshal(raw, &orgInterface)

		if err := ident.UpdateOrganization(*token.AccessToken, *common.Organization.ID, orgInterface); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		result, _ := json.MarshalIndent(secret, "", "\t")
		fmt.Printf("system secret: %s\n", string(result))
	} else {
		params["vault_id"] = vaults[0].ID
		system, err := baseline.CreateSystem(*token.AccessToken, common.WorkgroupID, params)
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		result, _ := json.MarshalIndent(system, "", "\t")
		fmt.Printf("system: %s\n", string(result))
	}
}

func systemTypePrompt() {
	systemTypes := make([]string, 2)
	systemTypes[0] = sapSystemIdentifier
	systemTypes[1] = servicenowSystemIdentifier

	prompt := promptui.Select{
		Label: "System Type",
		Items: systemTypes,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("failed to initialize system; %s", err.Error())
		os.Exit(1)
	}

	systemType = systemTypes[i]
}

func systemPrompt(params *map[string]interface{}, token string) {
	middlewareTypes := make([]string, 4)
	middlewareTypes[0] = systemNoMiddlewareIdentifier
	middlewareTypes[1] = systemInboundOnlyMiddlewareIdentifier
	middlewareTypes[2] = systemOutboundOnlyMiddlewareIdentifier
	middlewareTypes[3] = systemInboundAndOutboundMiddlewareIdentifier

	if systemMiddlewareType == "" {
		prompt := promptui.Select{
			Label: "Middleware Type",
			Items: middlewareTypes,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		systemMiddlewareType = middlewareTypes[i]
	}

	systemNamePrompt()

	systemDescriptionPrompt()

	switch systemMiddlewareType {
	case systemNoMiddlewareIdentifier:
		systemNoMiddlewarePrompt()

		systemAuthMethodsPrompt(&noMiddlewareAuthMethod, &noMiddlewareUsername, &noMiddlewarePassword)

		systemClientCredentialsPrompt(&noMiddlewareRequireClientCredentials, &noMiddlewareClientID, &noMiddlewareClientSecret)

		systemAuth := map[string]interface{}{
			"method":                   noMiddlewareAuthMethod,
			"username":                 noMiddlewareUsername,
			"password":                 noMiddlewarePassword,
			"require_user_credentials": noMiddlewareRequireClientCredentials,
		}

		if noMiddlewareRequireClientCredentials {
			systemAuth["client_id"] = noMiddlewareClientID
			systemAuth["client_secret"] = noMiddlewareClientSecret
		}

		reachabilityParams := map[string]interface{}{
			"type":         systemType,
			"name":         systemName,
			"auth":         systemAuth,
			"endpoint_url": systemEndpointURL,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		system := map[string]interface{}{
			"type":         systemType,
			"name":         systemName,
			"auth":         systemAuth,
			"endpoint_url": systemEndpointURL,
		}

		if systemDescription != "" {
			system["description"] = systemDescription
		}

		*params = system
	case systemInboundOnlyMiddlewareIdentifier:
		systemInboundOnlyPrompt()

		systemAuthMethodsPrompt(&inboundAuthMethod, &inboundPassword, &inboundUsername)

		systemClientCredentialsPrompt(&inboundRequireClientCredentials, &inboundClientID, &inboundClientSecret)

		inboundMiddlewareAuth := map[string]interface{}{
			"method":                   inboundAuthMethod,
			"username":                 inboundUsername,
			"password":                 inboundPassword,
			"require_user_credentials": inboundRequireClientCredentials,
		}

		if inboundRequireClientCredentials {
			inboundMiddlewareAuth["client_id"] = inboundClientID
			inboundMiddlewareAuth["client_secret"] = inboundClientSecret
		}

		reachabilityParams := map[string]interface{}{
			"name":         systemName,
			"type":         systemType,
			"auth":         inboundMiddlewareAuth,
			"endpoint_url": systemInboundEndpointURL,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		system := map[string]interface{}{
			"type": systemType,
			"name": systemName,
			"middleware": map[string]interface{}{
				"inbound": map[string]interface{}{
					"name": systemInboundMiddleware,
					"url":  systemInboundEndpointURL,
					"auth": inboundMiddlewareAuth,
				},
			},
		}

		if systemDescription != "" {
			system["description"] = systemDescription
		}

		*params = system
	case systemOutboundOnlyMiddlewareIdentifier:
		systemOutboundOnlyPrompt()

		systemAuthMethodsPrompt(&outboundAuthMethod, &outboundUsername, &outboundPassword)

		systemClientCredentialsPrompt(&outboundRequireClientCredentials, &outboundClientID, &outboundClientSecret)

		outboundMiddlewareAuth := map[string]interface{}{
			"method":                   outboundAuthMethod,
			"username":                 outboundUsername,
			"password":                 outboundPassword,
			"require_user_credentials": outboundRequireClientCredentials,
		}

		if outboundRequireClientCredentials {
			outboundMiddlewareAuth["client_id"] = outboundClientID
			outboundMiddlewareAuth["client_secret"] = outboundClientSecret
		}

		reachabilityParams := map[string]interface{}{
			"type":         systemType,
			"name":         systemName,
			"auth":         outboundMiddlewareAuth,
			"endpoint_url": systemOutboundEndpointURL,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		system := map[string]interface{}{
			"type": systemType,
			"name": systemName,
			"middleware": map[string]interface{}{
				"outbound": map[string]interface{}{
					"name": systemOutboundMiddleware,
					"url":  systemOutboundEndpointURL,
					"auth": outboundMiddlewareAuth,
				},
			},
		}

		if systemDescription != "" {
			system["description"] = systemDescription
		}

		*params = system
	case systemInboundAndOutboundMiddlewareIdentifier:
		systemInboundOnlyPrompt()

		systemAuthMethodsPrompt(&inboundAuthMethod, &inboundPassword, &inboundUsername)

		systemClientCredentialsPrompt(&inboundRequireClientCredentials, &inboundClientID, &inboundClientSecret)

		inboundMiddlewareAuth := map[string]interface{}{
			"method":                   inboundAuthMethod,
			"username":                 inboundUsername,
			"password":                 inboundPassword,
			"require_user_credentials": inboundRequireClientCredentials,
		}

		if inboundRequireClientCredentials {
			inboundMiddlewareAuth["client_id"] = inboundClientID
			inboundMiddlewareAuth["client_secret"] = inboundClientSecret
		}

		reachabilityParams := map[string]interface{}{
			"type":         systemType,
			"name":         systemName,
			"auth":         inboundMiddlewareAuth,
			"endpoint_url": systemInboundEndpointURL,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		systemOutboundOnlyPrompt()

		systemAuthMethodsPrompt(&outboundAuthMethod, &outboundUsername, &outboundPassword)

		systemClientCredentialsPrompt(&outboundRequireClientCredentials, &outboundClientID, &outboundClientSecret)

		outboundMiddlewareAuth := map[string]interface{}{
			"method":                   outboundAuthMethod,
			"username":                 outboundUsername,
			"password":                 outboundPassword,
			"require_user_credentials": outboundRequireClientCredentials,
		}

		if outboundRequireClientCredentials {
			outboundMiddlewareAuth["client_id"] = outboundClientID
			outboundMiddlewareAuth["client_secret"] = outboundClientSecret
		}

		reachabilityParams = map[string]interface{}{
			"auth":         outboundMiddlewareAuth,
			"endpoint_url": systemOutboundEndpointURL,
			"name":         systemName,
			"type":         systemType,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		system := map[string]interface{}{
			"type": systemType,
			"name": systemName,
			"middleware": map[string]interface{}{
				"inbound": map[string]interface{}{
					"auth": inboundMiddlewareAuth,
					"name": systemInboundMiddleware,
					"url":  systemInboundEndpointURL,
				},
				"outbound": map[string]interface{}{
					"auth": outboundMiddlewareAuth,
					"name": systemOutboundMiddleware,
					"url":  systemOutboundEndpointURL,
				},
			},
		}

		if systemDescription != "" {
			system["description"] = systemDescription
		}

		*params = system
	default:
		fmt.Print("failed to initialize system; invalid middleware type")
	}
}

func systemNamePrompt() {
	if systemName == "" {
		prompt := promptui.Prompt{
			Label:    "Name",
			Validate: common.MandatoryValidation,
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		systemName = result
	}
}

func systemDescriptionPrompt() {
	if systemDescription == "" {
		prompt := promptui.Prompt{
			Label: "Description",
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		systemDescription = result
	}
}

func systemNoMiddlewarePrompt() {
	if systemEndpointURL == "" {
		prompt := promptui.Prompt{
			Label: "Endpoint URL",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("endpoint url is required")
				}

				if _, err := url.ParseRequestURI(s); err != nil {
					return fmt.Errorf("invalid url")
				}

				return nil
			},
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize baseline domain model; %s", err.Error())
			os.Exit(1)
		}

		systemEndpointURL = result
	}
}

func systemInboundOnlyPrompt() {
	middlewareOpts := make([]string, 2)
	middlewareOpts[0] = "Mulesoft"
	middlewareOpts[1] = "SAPPI"

	if systemInboundMiddleware == "" {
		prompt := promptui.Select{
			Label: "Inbound Middleware Type",
			Items: middlewareOpts,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		systemInboundMiddleware = middlewareOpts[i]
	} else {
		isValid := false
		for _, opt := range middlewareOpts {
			if systemInboundMiddleware == opt {
				isValid = true
			}
		}

		if !isValid {
			fmt.Print("failed to initialize system; invalid system inbound middleware type")
			os.Exit(1)
		}
	}

	if systemInboundEndpointURL == "" {
		prompt := promptui.Prompt{
			Label: "Inbound Middleware URL",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("inbound middleware url is required")
				}

				if _, err := url.ParseRequestURI(s); err != nil {
					return fmt.Errorf("invalid url")
				}

				return nil
			},
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		systemInboundEndpointURL = result
	} else {
		if _, err := url.ParseRequestURI(systemInboundEndpointURL); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
	}
}

func systemOutboundOnlyPrompt() {
	middlewareOpts := make([]string, 2)
	middlewareOpts[0] = "Mulesoft"
	middlewareOpts[1] = "SAPPI"

	if systemOutboundMiddleware == "" {
		prompt := promptui.Select{
			Label: "Outbound Middleware Type",
			Items: middlewareOpts,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		systemOutboundMiddleware = middlewareOpts[i]
	} else {
		isValid := false
		for _, opt := range middlewareOpts {
			if systemOutboundMiddleware == opt {
				isValid = true
			}
		}

		if !isValid {
			fmt.Print("failed to initialize system; invalid system outbound middleware type")
			os.Exit(1)
		}
	}

	if systemOutboundEndpointURL == "" {
		prompt := promptui.Prompt{
			Label: "Outbound Middleware URL",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("outbound middleware url is required")
				}

				if _, err := url.ParseRequestURI(s); err != nil {
					return fmt.Errorf("invalid url")
				}

				return nil
			},
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		systemOutboundEndpointURL = result
	} else {
		if _, err := url.ParseRequestURI(systemOutboundEndpointURL); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
	}
}

func systemAuthMethodsPrompt(method, username, password *string) {
	authMethods := make([]string, 1)
	authMethods[0] = basicAuthMethodIdentifier

	if *method == "" {
		prompt := promptui.Select{
			Label: "Authentication Method",
			Items: authMethods,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		*method = authMethods[i]
	} else {
		isValid := false
		for _, am := range authMethods {
			if *method == am {
				isValid = true
			}
		}

		if !isValid {
			fmt.Print("failed to initialize system; invalid system authentication method")
			os.Exit(1)
		}
	}

	if *username == "" {
		prompt := promptui.Prompt{
			Label:    "Username",
			Validate: common.MandatoryValidation,
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		*username = result
	}

	if *password == "" {
		prompt := promptui.Prompt{
			Label:    "Password",
			Validate: common.MandatoryValidation,
			Mask:     '*',
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		*password = result
	}
}

func systemClientCredentialsPrompt(requireCredentials *bool, clientID, clientSecret *string) {
	if !*requireCredentials {
		prompt := promptui.Prompt{
			IsConfirm: true,
			Label:     "Require Client Credentials",
		}

		_, err := prompt.Run()
		if err != nil {
			return
		}

		*requireCredentials = true
	}

	if *clientID == "" {
		prompt := promptui.Prompt{
			Label:    "Client ID",
			Validate: common.MandatoryValidation,
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		*clientID = result
	}

	if *clientSecret == "" {
		prompt := promptui.Prompt{
			Label:    "Client Secret",
			Validate: common.MandatoryValidation,
			Mask:     '*',
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		*clientSecret = result
	}
}

func init() {
	initBaselineSystemCmd.Flags().StringVar(&common.OrganizationID, "organization", os.Getenv("PROVIDE_ORGANIZATION_ID"), "organization identifier")
	initBaselineSystemCmd.Flags().StringVar(&common.WorkgroupID, "workgroup", "", "workgroup identifier")
	initBaselineSystemCmd.Flags().StringVar(&systemType, "system-type", "", "system type")

	initBaselineSystemCmd.Flags().StringVar(&systemName, "name", "", "name")
	initBaselineSystemCmd.Flags().StringVar(&systemDescription, "description", "", "description")

	initBaselineSystemCmd.Flags().StringVar(&systemMiddlewareType, "middleware-type", "", "system middleware type")

	initBaselineSystemCmd.Flags().StringVar(&systemEndpointURL, "middleware-endpoint", "", "system middleware endpoint url")

	initBaselineSystemCmd.Flags().StringVar(&systemInboundMiddleware, "inbound-middleware", "", "system inbound middleware type")
	initBaselineSystemCmd.Flags().StringVar(&systemInboundEndpointURL, "inbound-endpoint", "", "system inbound middleware endpoint url")

	initBaselineSystemCmd.Flags().StringVar(&systemOutboundMiddleware, "outbound-middleware", "", "system outbound middleware type")
	initBaselineSystemCmd.Flags().StringVar(&systemOutboundEndpointURL, "outbound-endpoint", "", "system outbound middleware endpoint url")

	initBaselineSystemCmd.Flags().StringVar(&noMiddlewareAuthMethod, "auth-method", "", "no middleware authentication method")
	initBaselineSystemCmd.Flags().StringVar(&noMiddlewareUsername, "auth-username", "", "no middleware authentication username")
	initBaselineSystemCmd.Flags().StringVar(&noMiddlewarePassword, "auth-password", "", "no middleware authentication password")
	initBaselineSystemCmd.Flags().BoolVar(&noMiddlewareRequireClientCredentials, "require-client-credentials", false, "require no middleware client credentials")
	initBaselineSystemCmd.Flags().StringVar(&noMiddlewareClientID, "client-id", "", "no middleware system client id")
	initBaselineSystemCmd.Flags().StringVar(&noMiddlewareClientSecret, "client-secret", "", "no middleware system client secret")

	initBaselineSystemCmd.Flags().StringVar(&inboundAuthMethod, "inbound-auth-method", "", "inbound middleware authentication method")
	initBaselineSystemCmd.Flags().StringVar(&inboundUsername, "inbound-auth-username", "", "inbound middleware authentication username")
	initBaselineSystemCmd.Flags().StringVar(&inboundPassword, "inbound-auth-password", "", "inbound middleware authentication password")
	initBaselineSystemCmd.Flags().BoolVar(&inboundRequireClientCredentials, "inbound-require-client-credentials", false, "require inbound middleware client credentials")
	initBaselineSystemCmd.Flags().StringVar(&inboundClientID, "inbound-client-id", "", "inbound middleware system client id")
	initBaselineSystemCmd.Flags().StringVar(&inboundClientSecret, "inbound-client-secret", "", "inbound middleware system client secret")

	initBaselineSystemCmd.Flags().StringVar(&outboundAuthMethod, "oubound-auth-method", "", "outbound authentication method")
	initBaselineSystemCmd.Flags().StringVar(&outboundUsername, "oubound-auth-username", "", "outbound authentication username")
	initBaselineSystemCmd.Flags().StringVar(&outboundPassword, "oubound-auth-password", "", "outbound authentication password")
	initBaselineSystemCmd.Flags().BoolVar(&outboundRequireClientCredentials, "oubound-require-client-credentials", false, "require outbound client credentials")
	initBaselineSystemCmd.Flags().StringVar(&outboundClientID, "oubound-client-id", "", "outbound system client id")
	initBaselineSystemCmd.Flags().StringVar(&outboundClientSecret, "oubound-client-secret", "", "outbound system client secret")

	initBaselineSystemCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
