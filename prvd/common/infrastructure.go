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

package common

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const InfrastructureTargetAWS = "aws"
const InfrastructureTargetAzure = "azure"

var (
	EngineID           string
	ProviderID         string
	Region             string
	TargetID           string
	Image              string
	HealthCheckPath    string
	TaskRole           string
	TCPIngressPorts    string
	UDPIngressPorts    string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AzureTenantID      string
	AzureClientID      string
	AzureClientSecret  string
)

func InfrastructureCredentialsConfigFactory() map[string]interface{} {
	var creds map[string]interface{}

	if TargetID == InfrastructureTargetAWS {
		accessKeyID, secretAccessKey := requireAWSCredentials()
		creds = map[string]interface{}{
			"aws_access_key_id":     accessKeyID,
			"aws_secret_access_key": secretAccessKey,
		}
	} else if TargetID == InfrastructureTargetAzure {
		tenantID, clientID, clientSecret, subscriptionID := requireAzureCredentials()
		creds = map[string]interface{}{
			"azure_tenant_id":       tenantID,
			"azure_client_id":       clientID,
			"azure_client_secret":   clientSecret,
			"azure_subscription_id": subscriptionID,
		}
	}

	return creds
}

func RequireInfrastructureFlags(cmd *cobra.Command, withImage bool) {
	cmd.Flags().StringVar(&TargetID, "target", "aws", "target infrastructure platform (i.e., aws or azure)")
	cmd.Flags().StringVar(&Region, "Region", "us-east-1", "target infrastructure Region")
	cmd.Flags().StringVar(&ProviderID, "provider", "docker", "infrastructure virtualization provider (i.e., docker)")
	if withImage {
		cmd.Flags().StringVar(&Image, "common.Image", "", "container common.Image name")
	}
}

func requireAWSCredentials() (string, string) {
	fmt.Print("AWS Access Key ID: ")
	reader := bufio.NewReader(os.Stdin)
	accessKeyID, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	accessKeyID = strings.Trim(accessKeyID, "\n")
	if accessKeyID == "" {
		log.Println("Failed to read AWS access key ID from stdin")
		os.Exit(1)
	}

	fmt.Print("AWS Secret Access Key: ")
	secretAccessKeyBytes, err := terminal.ReadPassword(0)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	secretAccessKey := strings.Trim(string(secretAccessKeyBytes[:]), "\n")
	if secretAccessKey == "" {
		log.Println("Failed to read AWS secret access key from stdin")
		os.Exit(1)
	}

	return accessKeyID, secretAccessKey
}

func requireAzureCredentials() (string, string, string, string) {
	fmt.Print("Azure Tenant ID: ")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Azure Subscription ID: ")
	reader = bufio.NewReader(os.Stdin)
	subscriptionID, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	subscriptionID = strings.Trim(subscriptionID, "\n")
	if subscriptionID == "" {
		log.Println("Failed to read Azure subscription ID from stdin")
		os.Exit(1)
	}

	tenantID, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	tenantID = strings.Trim(tenantID, "\n")
	if tenantID == "" {
		log.Println("Failed to read Azure tenant ID from stdin")
		os.Exit(1)
	}

	fmt.Print("Azure Client ID: ")
	reader = bufio.NewReader(os.Stdin)
	clientID, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	clientID = strings.Trim(clientID, "\n")
	if clientID == "" {
		log.Println("Failed to read Azure client ID from stdin")
		os.Exit(1)
	}

	fmt.Print("Azure Client Secret: ")
	clientSecretBytes, err := terminal.ReadPassword(0)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	clientSecret := strings.Trim(string(clientSecretBytes[:]), "\n")
	if clientSecret == "" {
		log.Println("Failed to read Azure client secret from stdin")
		os.Exit(1)
	}

	return tenantID, clientID, clientSecret, subscriptionID
}
