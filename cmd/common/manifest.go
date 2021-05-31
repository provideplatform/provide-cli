package common

import (
	"fmt"
	"log"
	"os"
)

var manifestSaveLabel = "Would you like to save the configuration above?"
var manifestSaveArgs = []string{"Yes", "No"}

type PortMapping struct {
	HostPort      int
	ContainerPort int
}

// Add all the commons
type Keys struct {
	APIEndpoint                string
	ApiContainerPort           int
	NatsContainerPort          int
	NatsWebsocketContainerPort int
	NatsStreamingContainerPort int
	PostgresContainerPort      int
	RedisContainerPort         int
	//PortMapping                PortMapping
	DockerNetworkID   string
	Optional          bool
	Name              string
	Port              int
	IdentPort         int
	NchainPort        int
	PrivacyPort       int
	VaultPort         int
	NatsPort          int
	NatsWebsocketPort int
	NatsStreamingPort int
	PostgresPort      int
	RedisPort         int

	ApiHostname                      string
	ConsumerHostname                 string
	IdentHostname                    string
	IdentConsumerHostname            string
	NchainHostname                   string
	NchainConsumerHostname           string
	NchainStatsdaemonHostname        string
	NchainReachabilitydaemonHostname string
	PrivacyHostname                  string
	PrivacyConsumerHostname          string
	VaultHostname                    string
	NatsHostname                     string
	NatsServerName                   string
	NatsStreamingHostname            string
	PostgresHostname                 string
	RedisHostname                    string
	RedisHosts                       string

	AutoRemove bool
	LogLevel   string

	BaselineOrganizationAddress string

	//  baselineOrganizationAPIEndpoint string
	BaselineRegistryContractAddress string
	BaselineWorkgroupID             string

	NchainBaselineNetworkID string

	JwtSignerPublicKey string
	NatsAuthToken      string

	IdentAPIHost   string
	IdentAPIScheme string

	NchainAPIHost   string
	NchainAPIScheme string

	WorkgroupAccessToken     string
	OrganizationRefreshToken string

	PrivacyAPIHost   string
	PrivacyAPIScheme string

	SorID  string
	SorURL string

	VaultAPIHost       string
	VaultAPIScheme     string
	VaultRefreshToken  string
	VaultSealUnsealKey string

	SapAPIHost     string
	SapAPIScheme   string
	SapAPIUsername string
	SapAPIPassword string
	SapAPIPath     string

	ServiceNowAPIHost     string
	ServiceNowAPIScheme   string
	ServiceNowAPIUsername string
	ServiceNowAPIPassword string
	ServiceNowAPIPath     string

	SalesforceAPIHost   string
	SalesforceAPIScheme string
	SalesforceAPIPath   string

	WithLocalVault   bool
	WithLocalIdent   bool
	WithLocalNChain  bool
	WithLocalPrivacy bool
}

func ManifestSave(content []byte) {
	if SelectInput(manifestSaveArgs, manifestSaveLabel) == "Yes" {
		path := FreeInput("Please specify the path where you would like to save the Manifest", "", NoValidation)
		name := FreeInput("Please specify the path where you would like to save the Manifest", "", NoValidation)

		f, err := os.Create(fmt.Sprintf("%s/%s.txt", path, name))

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		_, err2 := f.Write(content)

		if err2 != nil {
			log.Fatal(err2)
		}

		fmt.Printf("Manifest File Saved to %s", path)
	}
}

func manifestLoad(contents string) {
	//TODO ... Read the contents as a string marshal them into a JSON and return them.
}
