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

var sapInboundMiddleware string
var sapInboundEndpointURL string

var sapOutboundMiddleware string
var sapOutboundEndpointURL string

var sapAuthMethod string

var sapUsername string
var sapPassword string

var sapRequireClientCredentials bool

var sapClientID string
var sapClientSecret string

const sapSystemIdentifier = "sap"

const sapNoMiddlewareIdentifier = "No Middleware"
const sapInboundOnlyMiddlewareIdentifier = "Inbound Only"
const sapOutboundOnlyMiddlewareIdentifier = "Outbound Only"
const sapInboundAndOutboundMiddlewareIdentifier = "Inbound & Outbound"

const sapAuthMethodIdentifier = "Basic Auth"

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

	var value string

	switch strings.ToLower(systemType) {
	case sapSystemIdentifier:
		sapPrompt(&value)
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

	// TODO-- use baseline url validation method to validate endpoint url
	// TODO-- schemas

	common.AuthorizeOrganizationContext(true)

	token := common.RequireOrganizationToken()

	vaults, err := vault.ListVaults(token, map[string]interface{}{})
	if err != nil {
		fmt.Printf("failed to initialize system; %s", err.Error())
		os.Exit(1)
	}

	secret, err := vault.CreateSecret(token, vaults[0].ID.String(), params)
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

		if err := baseline.UpdateWorkgroup(token, common.Workgroup.ID.String(), wgInterface); err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}
	}

	common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs = append(common.Organization.Metadata.Workgroups[common.Workgroup.ID].SystemSecretIDs, &secret.ID)

	var orgInterface map[string]interface{}
	raw, _ := json.Marshal(common.Organization)
	json.Unmarshal(raw, &orgInterface)

	if err := ident.UpdateOrganization(token, *common.Organization.ID, orgInterface); err != nil {
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

func sapPrompt(params *string) {

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

		sapAuthMethodPrompt()

		sapClientCredentialsPrompt()

		systemAuth := map[string]interface{}{
			"method":                   sapAuthMethod,
			"username":                 sapUsername,
			"password":                 sapPassword,
			"require_user_credentials": sapRequireClientCredentials,
		}

		if sapRequireClientCredentials {
			systemAuth["client_id"] = sapClientID
			systemAuth["client_secret"] = sapClientSecret
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

		sapAuthMethodPrompt()

		sapClientCredentialsPrompt()

		inboundMiddlewareAuth := map[string]interface{}{
			"method":                   sapAuthMethod,
			"username":                 sapUsername,
			"password":                 sapPassword,
			"require_user_credentials": sapRequireClientCredentials,
		}

		if sapRequireClientCredentials {
			inboundMiddlewareAuth["client_id"] = sapClientID
			inboundMiddlewareAuth["client_secret"] = sapClientSecret
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

		sapAuthMethodPrompt()

		sapClientCredentialsPrompt()

		outboundMiddlewareAuth := map[string]interface{}{
			"method":                   sapAuthMethod,
			"username":                 sapUsername,
			"password":                 sapPassword,
			"require_user_credentials": sapRequireClientCredentials,
		}

		if sapRequireClientCredentials {
			outboundMiddlewareAuth["client_id"] = sapClientID
			outboundMiddlewareAuth["client_secret"] = sapClientSecret
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

		sapAuthMethodPrompt()

		sapClientCredentialsPrompt()

		inboundMiddlewareAuth := map[string]interface{}{
			"method":                   sapAuthMethod,
			"username":                 sapUsername,
			"password":                 sapPassword,
			"require_user_credentials": sapRequireClientCredentials,
		}

		if sapRequireClientCredentials {
			inboundMiddlewareAuth["client_id"] = sapClientID
			inboundMiddlewareAuth["client_secret"] = sapClientSecret
		}

		// HACK-- there is not currently a complete set of flags that supports inbound & outbound middleware; this is a shortcut to allow reusing prompts
		sapAuthMethod = ""
		sapUsername = ""
		sapPassword = ""
		sapUsername = ""
		sapRequireClientCredentials = false
		sapClientID = ""
		sapClientSecret = ""

		sapOutboundOnlyPrompt()

		sapAuthMethodPrompt()

		sapClientCredentialsPrompt()

		outboundMiddlewareAuth := map[string]interface{}{
			"method":                   sapAuthMethod,
			"username":                 sapUsername,
			"password":                 sapPassword,
			"require_user_credentials": sapRequireClientCredentials,
		}

		if sapRequireClientCredentials {
			outboundMiddlewareAuth["client_id"] = sapClientID
			outboundMiddlewareAuth["client_secret"] = sapClientSecret
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
			os.Exit(1)
			fmt.Print("failed to initialize system; invalid sap inbound middleware type")
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
			os.Exit(1)
			fmt.Print("failed to initialize system; invalid sap outbound middleware type")
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

func sapAuthMethodPrompt() {
	authMethods := make([]string, 1)
	authMethods[0] = sapAuthMethodIdentifier

	if sapAuthMethod == "" {
		prompt := promptui.Select{
			Label: "Authentication Method",
			Items: authMethods,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		sapAuthMethod = authMethods[i]
	} else {
		isValid := false
		for _, am := range authMethods {
			if sapAuthMethod == am {
				isValid = true
			}
		}

		if !isValid {
			os.Exit(1)
			fmt.Print("failed to initialize system; invalid sap authentication method")
		}
	}

	if sapUsername == "" {
		prompt := promptui.Prompt{
			Label:    "Username",
			Validate: common.MandatoryValidation,
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		sapUsername = result
	}

	if sapPassword == "" {
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

		sapPassword = result
	}
}

func sapClientCredentialsPrompt() {
	if !sapRequireClientCredentials {
		prompt := promptui.Prompt{
			IsConfirm: true,
			Label:     "Require Client Credentials",
		}

		_, err := prompt.Run()
		if err != nil {
			return
		}

		sapRequireClientCredentials = true
	}

	if sapClientID == "" {
		prompt := promptui.Prompt{
			Label:    "Client ID",
			Validate: common.MandatoryValidation,
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("failed to initialize system; %s", err.Error())
			os.Exit(1)
		}

		sapClientID = result
	}

	if sapClientSecret == "" {
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

		sapClientSecret = result
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

	initBaselineSystemCmd.Flags().StringVar(&sapAuthMethod, "auth-method", "", "sap authentication method")
	initBaselineSystemCmd.Flags().StringVar(&sapUsername, "auth-username", "", "sap authentication username")
	initBaselineSystemCmd.Flags().StringVar(&sapPassword, "auth-password", "", "sap authentication password")

	initBaselineSystemCmd.Flags().BoolVarP(&sapRequireClientCredentials, "require-client-credentials", "", false, "require sap client credentials")
	initBaselineSystemCmd.Flags().StringVar(&sapClientID, "client-id", "", "sap client id")
	initBaselineSystemCmd.Flags().StringVar(&sapClientSecret, "client-secret", "", "sap client secret")

	initBaselineSystemCmd.Flags().BoolVarP(&Optional, "optional", "", false, "List all the Optional flags")
}
