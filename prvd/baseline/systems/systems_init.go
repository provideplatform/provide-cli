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
	"strings"

	uuid "github.com/kthomas/go.uuid"
	"github.com/manifoldco/promptui"
	"github.com/provideplatform/provide-cli/prvd/common"
	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/api/vault"
	"github.com/spf13/cobra"
)

var systemType string

var sapDescription string

var sapMiddlewareType string

var sapEndpointURL string

// could combine the inbound and outbound sap vars into sapEndpointURL slices similar to sapAuthMethods, sapUsernames, etc
var sapInboundMiddleware string
var sapInboundEndpointURL string

var sapOutboundMiddleware string
var sapOutboundEndpointURL string

var sapAuthMethods []string

var sapUsernames []string
var sapPasswords []string

var sapRequireClientCredentials []bool

var sapClientIDs []string
var sapClientSecrets []string

const sapSystemIdentifier = "sap"

const sapNoMiddlewareIdentifier = "No Middleware"
const sapInboundOnlyMiddlewareIdentifier = "Inbound Only"
const sapOutboundOnlyMiddlewareIdentifier = "Outbound Only"
const sapInboundAndOutboundMiddlewareIdentifier = "Inbound & Outbound"

const sapAuthMethodsIdentifier = "Basic Auth"

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

	var value string

	switch strings.ToLower(systemType) {
	case sapSystemIdentifier:
		sapPrompt(&value, *token.AccessToken)
	default:
		fmt.Print("failed to initialize system; invalid system type")
		os.Exit(1)
	}

	params := map[string]interface{}{
		"description": sapDescription,
		"name":        "system",
		"type":        "system",
		"value":       value,
	}

	vaults, err := vault.ListVaults(*token.AccessToken, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to initialize system; %s", err.Error())
		os.Exit(1)
	}

	secret, err := vault.CreateSecret(*token.AccessToken, vaults[0].ID.String(), params)
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
	raw, _ := json.Marshal(common.Organization)
	json.Unmarshal(raw, &orgInterface)

	if err := ident.UpdateOrganization(*token.AccessToken, *common.Organization.ID, orgInterface); err != nil {
		fmt.Printf("failed to initialize system; %s", err.Error())
		os.Exit(1)
	}

	result, _ := json.MarshalIndent(secret, "", "\t")
	fmt.Printf("%s\n", string(result))
}

func systemTypePrompt() {
	systemTypes := make([]string, 1)
	systemTypes[0] = sapSystemIdentifier

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

func sapPrompt(params *string, token string) {
	middlewareTypes := make([]string, 4)
	middlewareTypes[0] = sapNoMiddlewareIdentifier
	middlewareTypes[1] = sapInboundOnlyMiddlewareIdentifier
	middlewareTypes[2] = sapOutboundOnlyMiddlewareIdentifier
	middlewareTypes[3] = sapInboundAndOutboundMiddlewareIdentifier

	if sapMiddlewareType == "" {
		prompt := promptui.Select{
			Label: "Middleware Type",
			Items: middlewareTypes,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		sapMiddlewareType = middlewareTypes[i]
	}

	switch sapMiddlewareType {
	case sapNoMiddlewareIdentifier:
		sapDescriptionPrompt()

		sapNoMiddlewarePrompt()

		sapAuthMethodsPrompt(&sapAuthMethods[0], &sapUsernames[0], &sapPasswords[0])

		sapClientCredentialsPrompt(&sapRequireClientCredentials[0], &sapClientIDs[0], &sapClientSecrets[0])

		systemAuth := map[string]interface{}{
			"method":                   sapAuthMethods[0],
			"username":                 sapUsernames[0],
			"password":                 sapPasswords[0],
			"require_user_credentials": sapRequireClientCredentials[0],
		}

		if sapRequireClientCredentials[0] {
			systemAuth["client_id"] = sapClientIDs[0]
			systemAuth["client_secret"] = sapClientSecrets[0]
		}

		reachabilityParams := map[string]interface{}{
			"auth":         systemAuth,
			"endpoint_url": sapEndpointURL,
			"name":         "system",
			"type":         systemType,
			"description":  sapDescription,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		system := map[string]interface{}{
			"name":        "system",
			"description": sapDescription,
			"type":        "system",
			"value": map[string]interface{}{
				"auth":         systemAuth,
				"endpoint_url": sapEndpointURL,
				"name":         sapSystemIdentifier,
				"system":       strings.ToUpper(sapSystemIdentifier),
				"type":         sapNoMiddlewareIdentifier,
			},
		}

		value, err := json.Marshal(system)
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
		*params = string(value)
	case sapInboundOnlyMiddlewareIdentifier:
		sapDescriptionPrompt()

		sapInboundOnlyPrompt()

		sapAuthMethodsPrompt(&sapAuthMethods[0], &sapUsernames[0], &sapPasswords[0])

		sapClientCredentialsPrompt(&sapRequireClientCredentials[0], &sapClientIDs[0], &sapClientSecrets[0])

		inboundMiddlewareAuth := map[string]interface{}{
			"method":                   sapAuthMethods[0],
			"username":                 sapUsernames[0],
			"password":                 sapPasswords[0],
			"require_user_credentials": sapRequireClientCredentials[0],
		}

		if sapRequireClientCredentials[0] {
			inboundMiddlewareAuth["client_id"] = sapClientIDs[0]
			inboundMiddlewareAuth["client_secret"] = sapClientSecrets[0]
		}

		reachabilityParams := map[string]interface{}{
			"auth":         inboundMiddlewareAuth,
			"endpoint_url": sapInboundEndpointURL,
			"name":         "system",
			"type":         systemType,
			"description":  sapDescription,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		system := map[string]interface{}{
			"name":        "system",
			"description": sapDescription,
			"type":        "system",
			"value": map[string]interface{}{
				"name":   sapSystemIdentifier,
				"system": strings.ToUpper(sapSystemIdentifier),
				"type":   sapInboundOnlyMiddlewareIdentifier,
				"middleware": map[string]interface{}{
					"inbound": map[string]interface{}{
						"auth": inboundMiddlewareAuth,
						"name": sapInboundMiddleware,
						"url":  sapInboundEndpointURL,
					},
				},
			},
		}

		value, err := json.Marshal(system)
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
		*params = string(value)
	case sapOutboundOnlyMiddlewareIdentifier:
		sapDescriptionPrompt()

		sapOutboundOnlyPrompt()

		sapAuthMethodsPrompt(&sapAuthMethods[0], &sapUsernames[0], &sapPasswords[0])

		sapClientCredentialsPrompt(&sapRequireClientCredentials[0], &sapClientIDs[0], &sapClientSecrets[0])

		outboundMiddlewareAuth := map[string]interface{}{
			"method":                   sapAuthMethods[0],
			"username":                 sapUsernames[0],
			"password":                 sapPasswords[0],
			"require_user_credentials": sapRequireClientCredentials[0],
		}

		if sapRequireClientCredentials[0] {
			outboundMiddlewareAuth["client_id"] = sapClientIDs[0]
			outboundMiddlewareAuth["client_secret"] = sapClientSecrets[0]
		}

		reachabilityParams := map[string]interface{}{
			"auth":         outboundMiddlewareAuth,
			"endpoint_url": sapOutboundEndpointURL,
			"name":         "system",
			"type":         systemType,
			"description":  sapDescription,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		system := map[string]interface{}{
			"name":        "system",
			"description": sapDescription,
			"type":        "system",
			"value": map[string]interface{}{
				"name":   sapSystemIdentifier,
				"system": strings.ToUpper(sapSystemIdentifier),
				"type":   sapOutboundOnlyMiddlewareIdentifier,
				"middleware": map[string]interface{}{
					"outbound": map[string]interface{}{
						"auth": outboundMiddlewareAuth,
						"name": sapOutboundMiddleware,
						"url":  sapOutboundEndpointURL,
					},
				},
			},
		}

		value, err := json.Marshal(system)
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
		*params = string(value)
	case sapInboundAndOutboundMiddlewareIdentifier:
		sapDescriptionPrompt()

		sapInboundOnlyPrompt()

		sapAuthMethodsPrompt(&sapAuthMethods[0], &sapUsernames[0], &sapPasswords[0])

		sapClientCredentialsPrompt(&sapRequireClientCredentials[0], &sapClientIDs[0], &sapClientSecrets[0])

		inboundMiddlewareAuth := map[string]interface{}{
			"method":                   sapAuthMethods[0],
			"username":                 sapUsernames[0],
			"password":                 sapPasswords[0],
			"require_user_credentials": sapRequireClientCredentials[0],
		}

		if sapRequireClientCredentials[0] {
			inboundMiddlewareAuth["client_id"] = sapClientIDs[0]
			inboundMiddlewareAuth["client_secret"] = sapClientSecrets[0]
		}

		reachabilityParams := map[string]interface{}{
			"auth":         inboundMiddlewareAuth,
			"endpoint_url": sapInboundEndpointURL,
			"name":         "system",
			"type":         systemType,
			"description":  sapDescription,
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		sapOutboundOnlyPrompt()

		sapAuthMethodsPrompt(&sapAuthMethods[1], &sapUsernames[1], &sapPasswords[1])

		sapClientCredentialsPrompt(&sapRequireClientCredentials[1], &sapClientIDs[1], &sapClientSecrets[1])

		outboundMiddlewareAuth := map[string]interface{}{
			"method":                   sapAuthMethods[1],
			"username":                 sapUsernames[1],
			"password":                 sapPasswords[1],
			"require_user_credentials": sapRequireClientCredentials[1],
		}

		if sapRequireClientCredentials[1] {
			outboundMiddlewareAuth["client_id"] = sapClientIDs[1]
			outboundMiddlewareAuth["client_secret"] = sapClientSecrets[1]
		}

		reachabilityParams = map[string]interface{}{
			"subject_account_id": common.SHA256(fmt.Sprintf("%s.%s", common.OrganizationID, common.WorkgroupID)),
			"system": map[string]interface{}{
				"auth":         outboundMiddlewareAuth,
				"endpoint_url": sapOutboundEndpointURL,
				"name":         "system",
				"type":         systemType,
				"description":  sapDescription,
			},
		}

		if err := baseline.SystemReachability(token, reachabilityParams); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		system := map[string]interface{}{
			"name":        "system",
			"description": sapDescription,
			"type":        "system",
			"value": map[string]interface{}{
				"name":   sapSystemIdentifier,
				"system": strings.ToUpper(sapSystemIdentifier),
				"type":   sapInboundAndOutboundMiddlewareIdentifier,
				"middleware": map[string]interface{}{
					"inbound": map[string]interface{}{
						"auth": inboundMiddlewareAuth,
						"name": sapInboundMiddleware,
						"url":  sapInboundEndpointURL,
					},
					"outbound": map[string]interface{}{
						"auth": outboundMiddlewareAuth,
						"name": sapOutboundMiddleware,
						"url":  sapOutboundEndpointURL,
					},
				},
			},
		}

		value, err := json.Marshal(system)
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
		*params = string(value)
	default:
		fmt.Print("failed to initialize system; invalid middleware type")
	}
}

func sapDescriptionPrompt() {
	if sapDescription == "" {
		prompt := promptui.Prompt{
			Label:    "System Description",
			Validate: common.MandatoryValidation,
		}

		result, err := prompt.Run()
		if err != nil {
			os.Exit(1)
			return
		}

		sapDescription = result
	}
}

func sapNoMiddlewarePrompt() {
	if sapEndpointURL == "" {
		prompt := promptui.Prompt{
			Label: "SAP Endpoint URL",
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

		sapEndpointURL = result
	}
}

func sapInboundOnlyPrompt() {
	middlewareOpts := make([]string, 2)
	middlewareOpts[0] = "Mulesoft"
	middlewareOpts[1] = "SAPPI"

	if sapInboundMiddleware == "" {
		prompt := promptui.Select{
			Label: "SAP Inbound Middleware Type",
			Items: middlewareOpts,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		sapInboundMiddleware = middlewareOpts[i]
	} else {
		isValid := false
		for _, opt := range middlewareOpts {
			if sapInboundMiddleware == opt {
				isValid = true
			}
		}

		if !isValid {
			fmt.Print("failed to initialize system; invalid sap inbound middleware type")
			os.Exit(1)
		}
	}

	if sapInboundEndpointURL == "" {
		prompt := promptui.Prompt{
			Label: "SAP Inbound Middleware URL",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("sap inbound middleware url is required")
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

		sapInboundEndpointURL = result
	} else {
		if _, err := url.ParseRequestURI(sapInboundEndpointURL); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
	}
}

func sapOutboundOnlyPrompt() {
	middlewareOpts := make([]string, 2)
	middlewareOpts[0] = "Mulesoft"
	middlewareOpts[1] = "SAPPI"

	if sapOutboundMiddleware == "" {
		prompt := promptui.Select{
			Label: "SAP Outbound Middleware Type",
			Items: middlewareOpts,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		sapOutboundMiddleware = middlewareOpts[i]
	} else {
		isValid := false
		for _, opt := range middlewareOpts {
			if sapOutboundMiddleware == opt {
				isValid = true
			}
		}

		if !isValid {
			fmt.Print("failed to initialize system; invalid sap outbound middleware type")
			os.Exit(1)
		}
	}

	if sapOutboundEndpointURL == "" {
		prompt := promptui.Prompt{
			Label: "SAP Outbound Middleware URL",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("sap outbound middleware url is required")
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

		sapOutboundEndpointURL = result
	} else {
		if _, err := url.ParseRequestURI(sapOutboundEndpointURL); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
	}
}

func sapAuthMethodsPrompt(method, username, password *string) {
	authMethods := make([]string, 1)
	authMethods[0] = sapAuthMethodsIdentifier

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
			fmt.Print("failed to initialize system; invalid sap authentication method")
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

func sapClientCredentialsPrompt(requireCredentials *bool, clientID, clientSecret *string) {
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

	initBaselineSystemCmd.Flags().StringVar(&sapDescription, "description", "", "description")

	initBaselineSystemCmd.Flags().StringVar(&sapMiddlewareType, "middleware-type", "", "sap middleware type")

	initBaselineSystemCmd.Flags().StringVar(&sapEndpointURL, "middleware-endpoint", "", "sap middleware endpoint url")

	initBaselineSystemCmd.Flags().StringVar(&sapInboundMiddleware, "inbound-middleware", "", "sap inbound middleware type")
	initBaselineSystemCmd.Flags().StringVar(&sapInboundEndpointURL, "inbound-endpoint", "", "sap inbound middleware endpoint url")

	initBaselineSystemCmd.Flags().StringVar(&sapOutboundMiddleware, "outbound-middleware", "", "sap outbound middleware type")
	initBaselineSystemCmd.Flags().StringVar(&sapOutboundEndpointURL, "outbound-endpoint", "", "sap outbound middleware endpoint url")

	initBaselineSystemCmd.Flags().StringArrayVar(&sapAuthMethods, "auth-methods", []string{"", ""}, "sap authentication methods")
	initBaselineSystemCmd.Flags().StringArrayVar(&sapUsernames, "auth-usernames", []string{"", ""}, "sap authentication usernames")
	initBaselineSystemCmd.Flags().StringArrayVar(&sapPasswords, "auth-passwords", []string{"", ""}, "sap authentication passwords")

	initBaselineSystemCmd.Flags().BoolSliceVar(&sapRequireClientCredentials, "require-client-credentials", []bool{false, false}, "require sap client credentials")
	initBaselineSystemCmd.Flags().StringArrayVar(&sapClientIDs, "client-ids", []string{"", ""}, "sap client ids")
	initBaselineSystemCmd.Flags().StringArrayVar(&sapClientSecrets, "client-secrets", []string{"", ""}, "sap client secrets")

	initBaselineSystemCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
