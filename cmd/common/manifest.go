package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var manifestSaveLabel = "Would you like to save the configuration above?"
var manifestLoadLabel = "Would you like to load the configuration from a file?"
var selectInputArgs = []string{"Yes", "No"}
var LoadedKeys Keys

type PortMapping struct {
	HostPort      int
	ContainerPort int
}

// Add all the commons
type Keys struct {
	APIEndpoint                      string `json:"apiEndpoint"`
	MessagingEndpoint                string `json:"messagingEndpoint"`
	Tunnel                           bool   `json:"tunnel"`
	ExposeAPITunnel                  bool   `json:"exposeAPITunnel"`
	ExposeMessagingTunnel            bool   `json:"exposeMessagingTunnel"`
	ApiContainerPort                 int    `json:"apiContainerPort"`
	NatsContainerPort                int    `json:"natsContainerPort"`
	NatsWebsocketContainerPort       int    `json:"natsWebsocketContainerPort"`
	NatsStreamingContainerPort       int    `json:"natsStreamingContainerPort"`
	PostgresContainerPort            int    `json:"postgresContainerPort"`
	RedisContainerPort               int    `json:"redisContainerPort"`
	DockerNetworkID                  string `json:"dockerNetworkID"`
	Name                             string `json:"name"`
	Port                             int    `json:"port"`
	IdentPort                        int    `json:"identPort"`
	NchainPort                       int    `json:"nchainPort"`
	PrivacyPort                      int    `json:"privacyPort"`
	VaultPort                        int    `json:"vaultPort"`
	NatsPort                         int    `json:"natsPort"`
	NatsWebsocketPort                int    `json:"natsWebsocketPort"`
	NatsStreamingPort                int    `json:"natsStreamingPort"`
	PostgresPort                     int    `json:"postgresPort"`
	RedisPort                        int    `json:"reddisPort"`
	ApiHostname                      string `json:"apiHostName"`
	ConsumerHostname                 string `json:"consumerHostName"`
	IdentHostname                    string `json:"identHostName"`
	IdentConsumerHostname            string `json:"idemtConsumerHostName"`
	NchainHostname                   string `json:"nchainHostName"`
	NchainConsumerHostname           string `json:"nchainConsumerHostName"`
	NchainStatsdaemonHostname        string `json:"nchainStatsdaemonHostname"`
	NchainReachabilitydaemonHostname string `json:"nchainReachabilitydaemonHostname"`
	PrivacyHostname                  string `json:"privacyHostname"`
	PrivacyConsumerHostname          string `json:"privacyConsumerHostname"`
	VaultHostname                    string `json:"vaultHostName"`
	NatsHostname                     string `json:"natsHostName"`
	NatsServerName                   string `json:"natsServerName"`
	NatsStreamingHostname            string `json:"natsStreamingHostName"`
	PostgresHostname                 string `json:"postgresHostName"`
	RedisHostname                    string `json:"redisHostName"`
	RedisHosts                       string `json:"redisHosts"`
	AutoRemove                       bool   `json:"autoRemove"`
	LogLevel                         string `json:"logLevel"`
	BaselineOrganizationAddress      string `json:"baselineOrganizationAddress"`
	BaselineRegistryContractAddress  string `json:"baselineOrganizationContractAddress"`
	BaselineWorkgroupID              string `json:"baselineWorkgroupID"`
	NchainBaselineNetworkID          string `json:"nchainBaselineNetworkID"`
	JwtSignerPublicKey               string `json:"jwtSignerPublicKey"`
	NatsAuthToken                    string `json:"natsAuthToken"`
	IdentAPIHost                     string `json:"identAPIHost"`
	IdentAPIScheme                   string `json:"identAPIScheme"`
	NchainAPIHost                    string `json:"nchainAPIHost"`
	NchainAPIScheme                  string `json:"nchainAPIScheme"`
	WorkgroupAccessToken             string `json:"workgroupAccessToken"`
	OrganizationRefreshToken         string `json:"organizationRefreshToken"`
	PrivacyAPIHost                   string `json:"privacyAPIHost"`
	PrivacyAPIScheme                 string `json:"privacyAPIScheme"`
	SorID                            string `json:"SorID"`
	SorURL                           string `json:"SorURL"`
	VaultAPIHost                     string `json:"vaultAPIHost"`
	VaultAPIScheme                   string `json:"vaultAPIScheme"`
	VaultRefreshToken                string `json:"vaultRefreshToken"`
	VaultSealUnsealKey               string `json:"vaultUnsealToken"`
	SapAPIHost                       string `json:"sapAPIHost"`
	SapAPIScheme                     string `json:"sapAPIScheme"`
	SapAPIUsername                   string `json:"sapAPIUsername"`
	SapAPIPassword                   string `json:"sapAPIPassword"`
	SapAPIPath                       string `json:"sapAPIPath"`
	ServiceNowAPIHost                string `json:"servicenowAPIHost"`
	ServiceNowAPIScheme              string `json:"servicenowAPIScheme"`
	ServiceNowAPIUsername            string `json:"servicenowAPIUserName"`
	ServiceNowAPIPassword            string `json:"servicenowAPIPassword"`
	ServiceNowAPIPath                string `json:"servicenowAPIPath"`
	SalesforceAPIHost                string `json:"salesforceAPIHost"`
	SalesforceAPIScheme              string `json:"salesforceAPIScheme"`
	SalesforceAPIPath                string `json:"salesforceAPIPath"`
	WithLocalVault                   bool   `json:"withLocalVault"`
	WithLocalIdent                   bool   `json:"withLocalIdent"`
	WithLocalNChain                  bool   `json:"withLocalNchain"`
	WithLocalPrivacy                 bool   `json:"withLocalPrivacy"`
	//  PortMapping                PortMapping
	//  baselineOrganizationAPIEndpoint string
}

func ManifestSave(content []byte) {
	if SelectInput(selectInputArgs, manifestSaveLabel) == "Yes" {
		path := FreeInput("Please specify the path where you would like to save the Manifest", "", NoValidation)
		name := FreeInput("Please specify the path where you would like to save the Manifest", "", NoValidation)

		f, err := os.Create(fmt.Sprintf("%s/%s.json", path, name))

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

func ManifestLoad() bool {
	if SelectInput(selectInputArgs, manifestLoadLabel) == "Yes" {
		path := FreeInput("Please specify the path where you would like to load the configuration from", "", MandatoryValidation)

		data, _ := ioutil.ReadFile(path)
		json.Unmarshal(data, &LoadedKeys)
		fmt.Println(LoadedKeys.ApiHostname)

		return true
	}
	return false
	//TODO ... Read the contents as a string marshal them into a JSON and return them.

}
