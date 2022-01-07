package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/provideplatform/provide-go/api/ident"
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
	AccessTokenConfigKey            = "access-token"      // user-scoped API access token key
	RefreshTokenConfigKey           = "refresh-token"     // user-scoped API refresh token key
	APIAccessTokenConfigKeyPartial  = "api-token"         // app- or org-scoped API token key
	APIRefreshTokenConfigKeyPartial = "api-refresh-token" // app- or org-scoped API token key
	AccountConfigKeyPartial         = "account"           // app-scoped account ID key
	OrganizationConfigKeyPartial    = "organization"      // app-scoped organization ID key
	WalletConfigKeyPartial          = "wallet"            // app-scoped HD wallet ID key
	UserInfoConfigKey               = "prvd-user-info"    // details of the currently auth'd user
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
				if err != nil {
					fmt.Printf("WARNING: failed to write configuration; %s", err.Error())
				}
			}
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("WARNING: failed to read configuration; %s", err.Error())
	} else {
		os.Chmod(viper.ConfigFileUsed(), 0600)

		if Verbose {
			fmt.Println("Using configuration:", viper.ConfigFileUsed())
		}
	}
}

func StoreUserDetails(user *ident.User) error {
	json, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error storing user details; %s", err)
		return err
	}

	viper.Set(UserInfoConfigKey, string(json))
	viper.WriteConfig()

	return nil
}

func GetUserDetails() (*ident.User, error) {
	if viper.IsSet(UserInfoConfigKey) {
		raw := viper.GetString(UserInfoConfigKey)
		var user ident.User
		err := json.Unmarshal([]byte(raw), &user)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, fmt.Errorf("please authenticate to retrieve your user info")
}

func RequireUserAccessToken() string {
	token := ""
	if viper.IsSet(AccessTokenConfigKey) {
		token = viper.GetString(AccessTokenConfigKey)
	}

	if token == "" {
		log.Printf("Authorized API access token required in prvd configuration; run 'authenticate'")
		os.Exit(1)
	}

	if isTokenExpired(token) {
		refreshToken()
		token = viper.GetString(AccessTokenConfigKey)
	}

	return token
}

func refreshToken() {
	refreshToken := ""
	if viper.IsSet(RefreshTokenConfigKey) {
		refreshToken = viper.GetString(RefreshTokenConfigKey)
	}

	resp, err := ident.CreateToken(refreshToken, map[string]interface{}{
		"grant_type": "refresh_token",
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if resp != nil {
		CacheAccessRefreshToken(resp)
	}
}

func CacheAccessRefreshToken(token *ident.Token) {
	if token.AccessToken != nil {
		viper.Set(AccessTokenConfigKey, *token.AccessToken)
	}

	if token.RefreshToken != nil {
		viper.Set(RefreshTokenConfigKey, *token.RefreshToken)
	}

	viper.WriteConfig()
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
		token = RequireUserAccessToken()
	}

	if token == "" {
		log.Printf("Authorized API access token required in prvd configuration; run 'authenticate'")
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

func isTokenExpired(bearerToken string) bool {
	var jwtParser jwt.Parser
	token, _, err := jwtParser.ParseUnverified(bearerToken, jwt.MapClaims{})
	if err != nil {
		log.Printf("failed to parse JWT token on behalf of authorized user; %s", err.Error())
		os.Exit(1)
	}

	claims := token.Claims.(jwt.MapClaims)
	if exp, expOk := claims["exp"].(int64); expOk {
		expTime := time.Unix(exp, 0)
		now := time.Now()
		return expTime.Equal(now) || expTime.After(now)
	}

	return false
}
