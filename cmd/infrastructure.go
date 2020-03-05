package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const infrastructureTargetAWS = "aws"
const infrastructureTargetAzure = "azure"

var region string
var targetID string
var providerID string
var container string

var awsAccessKeyID string
var awsSecretAccessKey string

func infrastructureCredentialsConfigFactory() map[string]interface{} {
	var creds map[string]interface{}

	if targetID == infrastructureTargetAWS {
		accessKeyID, secretAccessKey := requireAWSCredentials()
		creds = map[string]interface{}{
			"aws_access_key_id":     accessKeyID,
			"aws_secret_access_key": secretAccessKey,
		}
	} else if targetID == infrastructureTargetAzure {
		tenantID, clientID, clientSecret := requireAzureCredentials()
		creds = map[string]interface{}{
			"azure_tenant_id":     tenantID,
			"azure_client_id":     clientID,
			"azure_client_secret": clientSecret,
		}
	}

	return creds
}

func requireInfrastructureFlags(cmd *cobra.Command, withContainer bool) {
	cmd.Flags().StringVar(&targetID, "target", "aws", "target infrastructure platform (i.e., aws or azure)")
	cmd.Flags().StringVar(&region, "region", "us-east-1", "target infrastructure region")
	cmd.Flags().StringVar(&providerID, "provider", "docker", "infrastructure virtualization provider (i.e., docker)")
	if withContainer {
		cmd.Flags().StringVar(&container, "container", "providenetwork-node", "infrastructure container (i.e., the name of the container image if using the docker provider)")
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
	reader = bufio.NewReader(os.Stdin)
	secretAccessKey, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	secretAccessKey = strings.Trim(secretAccessKey, "\n")
	if secretAccessKey == "" {
		log.Println("Failed to read AWS secret access key from stdin")
		os.Exit(1)
	}

	return accessKeyID, secretAccessKey
}

func requireAzureCredentials() (string, string, string) {
	fmt.Print("Azure Tenant ID: ")
	reader := bufio.NewReader(os.Stdin)
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
	reader = bufio.NewReader(os.Stdin)
	clientSecret, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	clientSecret = strings.Trim(clientSecret, "\n")
	if clientSecret == "" {
		log.Println("Failed to read Azure client secret from stdin")
		os.Exit(1)
	}

	return tenantID, clientID, clientSecret
}
