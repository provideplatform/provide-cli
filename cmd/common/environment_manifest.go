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
var LoadedManifest EnvManifest

type PortMapping struct {
	HostPort      int
	ContainerPort int
}

// Add all the commons
type EnvManifest struct {
	APIEndpoint                      string `json:"api_endpoint"`
	MessagingEndpoint                string `json:"messaging_endpoint"`
	Tunnel                           bool   `json:"tunnel_"`
	ExposeAPITunnel                  bool   `json:"expose_tunnel"`
	ExposeMessagingTunnel            bool   `json:"expose_messaging_tunnel"`
	ApiContainerPort                 int    `json:"api_port"`
	NatsContainerPort                int    `json:"nats_container_port"`
	NatsWebsocketContainerPort       int    `json:"nats_websocket_container_port"`
	NatsStreamingContainerPort       int    `json:"nats_streaming_container_port"`
	PostgresContainerPort            int    `json:"postgres_container_port"`
	RedisContainerPort               int    `json:"redis_port"`
	DockerNetworkID                  string `json:"docker_id"`
	Name                             string `json:"name"`
	Port                             int    `json:"port"`
	IdentPort                        int    `json:"ident_port"`
	NchainPort                       int    `json:"nchain_port"`
	PrivacyPort                      int    `json:"privacy_port"`
	VaultPort                        int    `json:"vault_port"`
	NatsPort                         int    `json:"nats_port"`
	NatsWebsocketPort                int    `json:"nats_websocket_port"`
	NatsStreamingPort                int    `json:"nats_streaming_port"`
	PostgresPort                     int    `json:"postgres_port"`
	RedisPort                        int    `json:"reddis_port"`
	ApiHostname                      string `json:"api_name"`
	ConsumerHostname                 string `json:"consumer_name"`
	IdentHostname                    string `json:"ident_name"`
	IdentConsumerHostname            string `json:"idemt_name"`
	NchainHostname                   string `json:"nchain_name"`
	NchainConsumerHostname           string `json:"nchain_consumer_hostname"`
	NchainStatsdaemonHostname        string `json:"nchain_hostname"`
	NchainReachabilitydaemonHostname string `json:"nchain_reachability_daemon_hostname"`
	PrivacyHostname                  string `json:"privacy_hostname"`
	PrivacyConsumerHostname          string `json:"privacy_consumer_hostname"`
	VaultHostname                    string `json:"vault_name"`
	NatsHostname                     string `json:"nats_name"`
	NatsServerName                   string `json:"nats_server_name"`
	NatsStreamingHostname            string `json:"nats_streaming_name"`
	PostgresHostname                 string `json:"postgres_name"`
	RedisHostname                    string `json:"redis_name"`
	RedisHosts                       string `json:"redis_hosts"`
	AutoRemove                       bool   `json:"auto_remove"`
	LogLevel                         string `json:"log_level"`
	BaselineOrganizationAddress      string `json:"baseline_address"`
	BaselineRegistryContractAddress  string `json:"baseline_registry_contractaddress"`
	BaselineWorkgroupID              string `json:"baseline_d"`
	NchainBaselineNetworkID          string `json:"nchain_d"`
	JwtSignerPublicKey               string `json:"jwt_key"`
	NatsAuthToken                    string `json:"nats_token"`
	IdentAPIHost                     string `json:"ident_host"`
	IdentAPIScheme                   string `json:"ident_scheme"`
	NchainAPIHost                    string `json:"nchain_host"`
	NchainAPIScheme                  string `json:"nchain_scheme"`
	WorkgroupAccessToken             string `json:"workgroup_token"`
	OrganizationRefreshToken         string `json:"organization_token"`
	PrivacyAPIHost                   string `json:"privacy_host"`
	PrivacyAPIScheme                 string `json:"privacy_scheme"`
	SorID                            string `json:"sor_d"`
	SorURL                           string `json:"sor_l"`
	VaultAPIHost                     string `json:"vault_host"`
	VaultAPIScheme                   string `json:"vault_scheme"`
	VaultRefreshToken                string `json:"vault_refresh_token"`
	VaultSealUnsealKey               string `json:"vault_un/seal_token"`
	SapAPIHost                       string `json:"sap_host"`
	SapAPIScheme                     string `json:"sap_scheme"`
	SapAPIUsername                   string `json:"sap_username"`
	SapAPIPassword                   string `json:"sap_password"`
	SapAPIPath                       string `json:"sap_path"`
	ServiceNowAPIHost                string `json:"servicenow_host"`
	ServiceNowAPIScheme              string `json:"servicenow_scheme"`
	ServiceNowAPIUsername            string `json:"servicenow_name"`
	ServiceNowAPIPassword            string `json:"servicenow_password"`
	ServiceNowAPIPath                string `json:"servicenow_path"`
	SalesforceAPIHost                string `json:"salesforce_host"`
	SalesforceAPIScheme              string `json:"salesforce_scheme"`
	SalesforceAPIPath                string `json:"salesforce_path"`
	WithLocalVault                   bool   `json:"with_vault"`
	WithLocalIdent                   bool   `json:"with_ident"`
	WithLocalNChain                  bool   `json:"with_nchain"`
	WithLocalPrivacy                 bool   `json:"with_privacy"`
	//  PortMapping                PortMapping
	//  baselineOrganizationAPIEndpoint string
}

func ManifestSave(content []byte) {
	if SelectInput(selectInputArgs, manifestSaveLabel) == "Yes" {
		path := FreeInput("Please specify the path where you would like to save the Manifest", "", NoValidation)
		f, err := os.Create(path)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		_, err2 := f.Write(content)

		if err2 != nil {
			log.Fatal(err2)
		}

		fmt.Printf("Manifest File Saved to %s \n", path)
		envProvideBaselineManifest = path
	}
}

func ManifestLoad() bool {
	fmt.Println(envProvideBaselineManifest)
	if envProvideBaselineManifest != "" {
		if SelectInput(selectInputArgs, manifestLoadLabel) == "Yes" {
			data, _ := ioutil.ReadFile(envProvideBaselineManifest)
			json.Unmarshal(data, &LoadedManifest)
			fmt.Println(LoadedManifest.ApiHostname)

			return true
		}
	}
	return false
}
