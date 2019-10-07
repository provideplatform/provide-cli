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
	}

	return creds
}

func requireInfrastructureFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&targetID, "target", "aws", "target infrastructure platform (i.e., aws or azure)")
	cmd.Flags().StringVar(&region, "region", "us-east-1", "infrastructure region")
	cmd.Flags().StringVar(&providerID, "provider", "docker", "infrastructure virtualization provider (i.e., docker)")
	cmd.Flags().StringVar(&container, "container", "providenetwork-node", "infrastructure container (i.e., the name of the container image if using the docker provider)")
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
