package cmd

import "github.com/spf13/cobra"

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
		creds = map[string]interface{}{
			"aws_access_key_id":     awsAccessKeyID,
			"aws_secret_access_key": awsSecretAccessKey,
		}
	}

	return creds
}
func requireInfrastructureFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&targetID, "target", "aws", "target infrastructure platform (i.e., aws or azure)")
	cmd.MarkFlagRequired("target")

	cmd.Flags().StringVar(&region, "region", "us-east-1", "infrastructure region")
	cmd.MarkFlagRequired("region")

	cmd.Flags().StringVar(&providerID, "provider", "docker", "infrastructure virtualization provider (i.e., docker)")
	cmd.MarkFlagRequired("provider")

	cmd.Flags().StringVar(&container, "container", "providenetwork-node", "infrastructure container (i.e., the name of the container image if using the docker provider)")
	cmd.MarkFlagRequired("container")

	if targetID == infrastructureTargetAWS {
		requireAWSCredentialsFlags(cmd)
	}
}

func requireAWSCredentialsFlags(cmd *cobra.Command) {
	// FIXME-- allow these to be prompted instead...

	cmd.Flags().StringVar(&awsAccessKeyID, "aws_access_key_id", "", "aws access key id")
	cmd.MarkFlagRequired("aws_access_key_id")

	cmd.Flags().StringVar(&awsSecretAccessKey, "aws_secret_access_key", "", "aws secret access key")
	cmd.MarkFlagRequired("aws_secret_access_key")
}
