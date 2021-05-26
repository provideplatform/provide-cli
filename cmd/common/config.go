package common

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	ASCIIBanner = `██████╗ ██████╗  ██████╗ ██╗   ██╗██╗██████╗ ███████╗
██╔══██╗██╔══██╗██╔═══██╗██║   ██║██║██╔══██╗██╔════╝
██████╔╝██████╔╝██║   ██║██║   ██║██║██║  ██║█████╗  
██╔═══╝ ██╔══██╗██║   ██║╚██╗ ██╔╝██║██║  ██║██╔══╝  
██║     ██║  ██║╚██████╔╝ ╚████╔╝ ██║██████╔╝███████╗
╚═╝     ╚═╝  ╚═╝ ╚═════╝   ╚═══╝  ╚═╝╚═════╝ ╚══════╝`

	// Viper downcases key names, so hyphenating for better readability.
	// 'Partial' keys are to be combined with the application ID they are associated with.
	// and NOT used by themselves.
	AuthTokenConfigKey              = "auth-token"        // user-scoped API token key
	APIAccessTokenConfigKeyPartial  = "api-token"         // app- or org-scoped API token key
	APIRefreshTokenConfigKeyPartial = "api-refresh-token" // app- or org-scoped API token key
	AccountConfigKeyPartial         = "account"           // app-scoped account ID key
	OrganizationConfigKeyPartial    = "organization"      // app-scoped organization ID key
	WalletConfigKeyPartial          = "wallet"            // app-scoped HD wallet ID key
)

var CfgFile string

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".provide-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".provide-cli")

		configPath := fmt.Sprintf("%s/.provide-cli.yaml", home)
		if err := viper.SafeWriteConfigAs(configPath); err != nil {
			if os.IsNotExist(err) {
				err = viper.WriteConfigAs(configPath)
			}
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		os.Chmod(viper.ConfigFileUsed(), 0600)

		if Verbose {
			fmt.Println("Using configuration:", viper.ConfigFileUsed())
		}
	}
}

func RequireUserAuthToken() string {
	token := ""
	if viper.IsSet(AuthTokenConfigKey) {
		token = viper.GetString(AuthTokenConfigKey)
	}

	if token == "" {
		log.Printf("Authorized API token required in prvd configuration; run 'authenticate'")
		os.Exit(1)
	}
	return token
}

func RequireApplicationToken() string {
	var token string
	tokenKey := BuildConfigKeyWithApp(APIAccessTokenConfigKeyPartial, ApplicationID)
	if viper.IsSet(tokenKey) {
		token = viper.GetString(tokenKey)
	}

	if token == "" {
		log.Printf("Authorized application API token required in prvd configuration; run 'prvd api_tokens init --application <id>'")
		os.Exit(1)
	}

	return token
}

func RequireOrganizationToken() string {
	var token string
	tokenKey := BuildConfigKeyWithOrg(APIAccessTokenConfigKeyPartial, OrganizationID)
	if viper.IsSet(tokenKey) {
		token = viper.GetString(tokenKey)
	}

	if token == "" {
		log.Printf("Authorized organization API token required in prvd configuration; run 'prvd api_tokens init --organization <id>'")
		os.Exit(1)
	}

	return token
}

func RequireAPIToken() string {
	var token string
	var appAPITokenKey string
	var orgAPITokenKey string
	if ApplicationID != "" {
		appAPITokenKey = BuildConfigKeyWithApp(APIAccessTokenConfigKeyPartial, ApplicationID)
	} else if OrganizationID != "" {
		orgAPITokenKey = BuildConfigKeyWithOrg(APIAccessTokenConfigKeyPartial, OrganizationID)
	}
	if viper.IsSet(appAPITokenKey) {
		token = viper.GetString(appAPITokenKey)
	} else if viper.IsSet(orgAPITokenKey) {
		token = viper.GetString(orgAPITokenKey)
	} else {
		token = RequireUserAuthToken()
	}

	if token == "" {
		log.Printf("Authorized API token required in prvd configuration; run 'authenticate'")
		os.Exit(1)
	}
	return token
}

// BuildConfigKeyWithApp combines the given key partial and app ID according to a consistent convention.
// Returns an empty string if the given appID is empty.
// Viper's getters likewise return empty strings when passed an empty string.
func BuildConfigKeyWithApp(keyPartial, appID string) string {
	if appID == "" {
		// Development-time debugging.
		log.Println("An application identifier is required for this operation")
		return ""
	}
	return fmt.Sprintf("%s.%s", appID, keyPartial)
}

// BuildConfigKeyWithOrg combines the given key partial and org ID according to a consistent convention.
// Returns an empty string if the given orgID is empty.
// Viper's getters likewise return empty strings when passed an empty string.
func BuildConfigKeyWithOrg(keyPartial, orgID string) string {
	if orgID == "" {
		// Development-time debugging.
		log.Println("An organization identifier is required for this operation")
		return ""
	}
	return fmt.Sprintf("%s.%s", orgID, keyPartial)
}

// BuildConfigKeyWithUser combines the given key partial and user ID according to a consistent convention.
// Returns an empty string if the given userID is empty.
// Viper's getters likewise return empty strings when passed an empty string.
func BuildConfigKeyWithUser(keyPartial, userID string) string {
	if userID == "" {
		// Development-time debugging.
		log.Println("A user identifier is required for this operation")
		return ""
	}
	return fmt.Sprintf("%s.%s", userID, keyPartial)
}
