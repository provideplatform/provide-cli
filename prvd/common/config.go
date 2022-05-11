package common

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/common/util"
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
	AccessTokenConfigKey         = "access-token"  // user-scoped API access token key
	RefreshTokenConfigKey        = "refresh-token" // user-scoped API refresh token key
	AccountConfigKeyPartial      = "account"       // app-scoped account ID key
	OrganizationConfigKeyPartial = "organization"  // app-scoped organization ID key
	WalletConfigKeyPartial       = "wallet"        // app-scoped HD wallet ID key
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

func RequireUserAccessToken() string {
	token := ""
	if viper.IsSet(AccessTokenConfigKey) {
		token = viper.GetString(AccessTokenConfigKey)
	}

	if token == "" || isTokenExpired(token) {
		log.Printf("Authorized API access token required in prvd configuration; run 'authenticate'")
		os.Exit(1)
	}

	if isTokenExpired(token) {
		refreshToken(token, nil)
		token = viper.GetString(AccessTokenConfigKey)
	}

	return token
}

func refreshToken(token string, id *string) {
	var refreshTokenKey string
	if id == nil {
		refreshTokenKey = RefreshTokenConfigKey
	} else {
		refreshTokenKey = BuildConfigKeyWithID(RefreshTokenConfigKey, *id)
	}

	var refreshToken string
	if viper.IsSet(refreshTokenKey) {
		refreshToken = viper.GetString(refreshTokenKey)
	}

	resp, err := ident.CreateToken(refreshToken, map[string]interface{}{
		"grant_type": "refresh_token",
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if resp != nil {
		CacheAccessRefreshToken(resp, id)
	}
}

func CacheAccessRefreshToken(token *ident.Token, id *string) {
	var accessTokenKey string
	var refreshTokenKey string

	if id == nil {
		accessTokenKey = AccessTokenConfigKey
		refreshTokenKey = RefreshTokenConfigKey
	} else {
		accessTokenKey = BuildConfigKeyWithID(AccessTokenConfigKey, *id)
		refreshTokenKey = BuildConfigKeyWithID(RefreshTokenConfigKey, *id)
	}

	if token.AccessToken != nil {
		viper.Set(accessTokenKey, *token.AccessToken)
	}

	if token.RefreshToken != nil {
		viper.Set(refreshTokenKey, *token.RefreshToken)
	}

	viper.WriteConfig()
}

func RequireApplicationToken() string {
	var token string
	tokenKey := BuildConfigKeyWithID(AccessTokenConfigKey, ApplicationID)
	if viper.IsSet(tokenKey) {
		token = viper.GetString(tokenKey)
	}

	if token == "" || isTokenExpired(token) {
		log.Printf("Authorized application API token required in prvd configuration; run 'prvd api_tokens init --application <id>'")
		os.Exit(1)
	}

	return token
}

func RequireOrganizationToken() string {
	var token string
	tokenKey := BuildConfigKeyWithID(AccessTokenConfigKey, OrganizationID)
	if viper.IsSet(tokenKey) {
		token = viper.GetString(tokenKey)
	}

	if token == "" || isTokenExpired(token) {
		err := RequireOrganization()
		if err == nil && PromptOrganizationAuthorization() {
			token = viper.GetString(tokenKey)
		}
	}

	return token
}

func PromptOrganizationAuthorization() bool {
	prompt := promptui.Prompt{
		IsConfirm: true,
		Label:     fmt.Sprintf("Authorize access/refresh token for %s", *Organization.Name),
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
		return false
	}

	if strings.ToLower(result) == "y" {
		return AuthorizeOrganization(true) == nil
	}

	return false
}

func AuthorizeOrganization(persist bool) error {
	token, err := ident.CreateToken(RequireUserAccessToken(), map[string]interface{}{
		"scope":           "offline_access",
		"organization_id": OrganizationID,
	})
	if err != nil {
		log.Printf("failed to authorize API access token on behalf of organization %s; %s", OrganizationID, err.Error())
		os.Exit(1)
	}

	if token.AccessToken != nil {
		OrganizationAccessToken = *token.AccessToken

		if token.RefreshToken != nil {
			OrganizationRefreshToken = *token.RefreshToken
		}

		if persist {
			// FIXME-- DRY this up (also exists in api_tokens_init.go)
			orgAPIAccessTokenKey := BuildConfigKeyWithID(AccessTokenConfigKey, OrganizationID)
			orgAPIRefreshTokenKey := BuildConfigKeyWithID(RefreshTokenConfigKey, OrganizationID)

			if token.AccessToken != nil {
				// fmt.Printf("Access token authorized for organization: %s\t%s\n", OrganizationID, *token.AccessToken)
				if !viper.IsSet(orgAPIAccessTokenKey) {
					viper.Set(orgAPIAccessTokenKey, *token.AccessToken)
					viper.WriteConfig()
				}
				if token.RefreshToken != nil {
					// fmt.Printf("Refresh token authorized for organization: %s\t%s\n", OrganizationID, *token.RefreshToken)
					if !viper.IsSet(orgAPIRefreshTokenKey) {
						viper.Set(orgAPIRefreshTokenKey, *token.RefreshToken)
						viper.WriteConfig()
					}
				}
			}
		}
	}

	return nil
}

func RequireAPIToken() string {
	var token string
	var tokenKey string
	var id *string

	if ApplicationID != "" {
		tokenKey = BuildConfigKeyWithID(AccessTokenConfigKey, ApplicationID)
		id = &ApplicationID
	} else if OrganizationID != "" {
		tokenKey = BuildConfigKeyWithID(AccessTokenConfigKey, OrganizationID)
		id = &OrganizationID
	} else {
		tokenKey = AccessTokenConfigKey
	}

	requireToken := func() string {
		if viper.IsSet(tokenKey) {
			token = viper.GetString(tokenKey)
			if isTokenExpired(token) {
				refreshToken(token, id)
				return viper.GetString(tokenKey)
			}
		}
		return RequireUserAccessToken()
	}

	token = requireToken()
	if token == "" {
		log.Printf("Authorized API access token required in prvd configuration; run 'authenticate'")
		os.Exit(1)
	}

	return token
}

// BuildConfigKeyWithID combines the given key partial and ID
// according to a consistent convention. Returns an empty string
// if the given id is empty. Viper's getters likewise return empty
// strings when passed an empty string.
func BuildConfigKeyWithID(keyPartial, id string) string {
	if id == "" {
		return ""
	}
	return fmt.Sprintf("%s.%s", id, keyPartial)
}

func isTokenExpired(bearerToken string) bool {
	token, _ := jwt.Parse(bearerToken, func(_jwtToken *jwt.Token) (interface{}, error) {
		// uncomment when enabling local verification
		var kid *string
		if kidhdr, ok := _jwtToken.Header["kid"].(string); ok {
			kid = &kidhdr
		}

		publicKey, _, _, _ := util.ResolveJWTKeypair(kid)
		if publicKey == nil {
			msg := "failed to resolve a valid JWT verification key"
			if kid != nil {
				msg = fmt.Sprintf("%s; invalid kid specified in header: %s", msg, *kid)
			} else {
				msg = fmt.Sprintf("%s; no default verification key configured", msg)
			}
			return nil, fmt.Errorf(msg)
		}

		return publicKey, nil
	})

	// TODO-- to enable this, enable caching of JWT keypairs locally so the above util.ResolveJWTKeypair(kid) successfully resolves
	// if err != nil {
	// 	return false
	// }

	claims := token.Claims.(jwt.MapClaims)
	if exp, expOk := claims["exp"].(float64); expOk {
		expTime := time.Unix(int64(exp), 0)
		now := time.Now()
		return expTime.Before(now) || expTime.Equal(now)
	}

	return false
}
