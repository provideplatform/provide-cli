package api_tokens

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/kthomas/go-pgputil"
	"github.com/provideservices/provide-cli/cmd/common"
	"github.com/provideservices/provide-go/api/ident"
	provide "github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/common/util"
	"golang.org/x/crypto/ssh"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jwtKeypairs map[string]*util.JWTKeypair

var scope string
var grantType string
var offlineAccess bool
var refreshToken bool
var optional bool

var apiTokensInitCmd = &cobra.Command{
	Use:   "init [--application 8fec625c-a8ad-4197-bb77-8b46d7aecd8f] [--organization 2209cf15-2402-4e25-b6b6-1c901b9dde69] [--offline-access] [--refresh-token]",
	Short: "Authorize a new API access or refresh token",
	Long:  `Authorize a new API token on behalf of the given application or organization`,
	Run:   createAPIToken,
}

// createAPIToken triggers the generation of an API token for the given network.
func createAPIToken(cmd *cobra.Command, args []string) {
	RequirePublicJWTVerifiers()

	userToken := common.RequireUserAccessToken()
	params := map[string]interface{}{}

	if scope != "" {
		params["scope"] = scope
	} else if offlineAccess {
		params["scope"] = "offline_access"
	}

	if grantType != "" {
		params["grant_type"] = grantType
	} else if refreshToken {
		params["grant_type"] = "refresh_token"
	}

	if common.ApplicationID != "" {
		token, err := provide.CreateApplicationToken(userToken, common.ApplicationID, params)
		if err != nil {
			log.Printf("Failed to authorize API token on behalf of application %s; %s", common.ApplicationID, err.Error())
			os.Exit(1)
		}

		appAPITokenKey := common.BuildConfigKeyWithApp(common.APIAccessTokenConfigKeyPartial, common.ApplicationID)
		appAPIRefreshTokenKey := common.BuildConfigKeyWithApp(common.APIRefreshTokenConfigKeyPartial, common.ApplicationID)
		var tkn string

		if token.Token != nil {
			fmt.Printf("API token authorized for application: %s\t%s\n", common.ApplicationID, *token.Token)
			tkn = *token.Token
		} else if token.AccessToken != nil {
			fmt.Printf("Access token authorized for application: %s\t%s\n", common.ApplicationID, *token.AccessToken)
			tkn = *token.AccessToken
		}

		if tkn != "" {
			if !viper.IsSet(appAPITokenKey) {
				viper.Set(appAPITokenKey, tkn)
				viper.WriteConfig()
			}

			if token.RefreshToken != nil {
				fmt.Printf("Refresh token authorized for application: %s\t%s\n", common.ApplicationID, *token.RefreshToken)
				if !viper.IsSet(appAPIRefreshTokenKey) {
					viper.Set(appAPIRefreshTokenKey, *token.RefreshToken)
					viper.WriteConfig()
				}
			}
		}
	} else if common.OrganizationID != "" {
		params["organization_id"] = common.OrganizationID
		token, err := provide.CreateToken(userToken, params)
		if err != nil {
			log.Printf("failed to authorize API access token on behalf of organization %s; %s", common.OrganizationID, err.Error())
			os.Exit(1)
		}

		orgAPIAccessTokenKey := common.BuildConfigKeyWithOrg(common.APIAccessTokenConfigKeyPartial, common.OrganizationID)
		orgAPIRefreshTokenKey := common.BuildConfigKeyWithOrg(common.APIRefreshTokenConfigKeyPartial, common.OrganizationID)

		if token.AccessToken != nil {
			fmt.Printf("Access token authorized for organization: %s\t%s\n", common.OrganizationID, *token.AccessToken)
			if !viper.IsSet(orgAPIAccessTokenKey) {
				viper.Set(orgAPIAccessTokenKey, *token.AccessToken)
				viper.WriteConfig()
			}
			if token.RefreshToken != nil {
				fmt.Printf("Refresh token authorized for organization: %s\t%s\n", common.OrganizationID, *token.RefreshToken)
				if !viper.IsSet(orgAPIRefreshTokenKey) {
					viper.Set(orgAPIRefreshTokenKey, *token.RefreshToken)
					viper.WriteConfig()
				}
			}
		} else {
			log.Printf("Failed to authorize API token on behalf of organization %s; no access/refresh pair returned", common.OrganizationID)
			os.Exit(1)
		}
	} else {
		// user token...
		token, err := provide.CreateToken(userToken, params)
		if err != nil {
			log.Printf("failed to authorize API access token on behalf of authorized user; %s", err.Error())
			os.Exit(1)
		}

		tkn, err := ParseJWT(userToken)
		if err != nil {
			log.Printf("failed to parse JWT token on behalf of authorized user; %s", err.Error())
			os.Exit(1)
		}
		claims, _ := tkn.Claims.(jwt.MapClaims)

		userID := strings.Split(claims["sub"].(string), ":")[1]
		userAPIAccessTokenKey := common.BuildConfigKeyWithUser(common.APIAccessTokenConfigKeyPartial, userID)
		userAPIRefreshTokenKey := common.BuildConfigKeyWithUser(common.APIRefreshTokenConfigKeyPartial, userID)

		if token.AccessToken != nil {
			fmt.Printf("Access token authorized for user: %s\t%s\n", common.OrganizationID, *token.AccessToken)
			if !viper.IsSet(userAPIAccessTokenKey) {
				viper.Set(userAPIAccessTokenKey, *token.AccessToken)
				viper.WriteConfig()
			}
			if token.RefreshToken != nil {
				fmt.Printf("Refresh token authorized for user: %s\t%s\n", common.OrganizationID, *token.RefreshToken)
				if !viper.IsSet(userAPIRefreshTokenKey) {
					viper.Set(userAPIRefreshTokenKey, *token.RefreshToken)
					viper.WriteConfig()
				}
			}
		} else {
			log.Printf("Failed to authorize API token on behalf of authorized user %s; no access/refresh pair returned", userID)
			os.Exit(1)
		}
	}
}

func RequirePublicJWTVerifiers() {
	jwtKeypairs = map[string]*util.JWTKeypair{}

	keys, err := ident.GetJWKs()
	if err != nil {
		log.Printf("failed to resolve ident jwt keys; %s", err.Error())
	} else {
		for _, key := range keys {
			publicKey, err := pgputil.DecodeRSAPublicKeyFromPEM([]byte(key.PublicKey))
			if err != nil {
				log.Printf("failed to parse ident JWT public key; %s", err.Error())
			}

			sshPublicKey, err := ssh.NewPublicKey(publicKey)
			if err != nil {
				log.Printf("failed to resolve JWT public key fingerprint; %s", err.Error())
			}
			fingerprint := ssh.FingerprintLegacyMD5(sshPublicKey)

			jwtKeypairs[fingerprint] = &util.JWTKeypair{
				Fingerprint:  fingerprint,
				PublicKey:    *publicKey,
				PublicKeyPEM: &key.PublicKey,
			}

			log.Printf("ident jwt public key configured for verification; fingerprint: %s", fingerprint)
		}
	}
}

func ParseJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(_jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := _jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("failed to resolve a valid JWT signing key; unsupported signing alg specified in header: %s", _jwtToken.Method.Alg())
		}

		var keypair *util.JWTKeypair

		var kid *string
		if kidhdr, ok := _jwtToken.Header["kid"].(string); ok {
			kid = &kidhdr
		}

		if kid != nil {
			keypair = jwtKeypairs[*kid]
		}

		if keypair == nil {
			for kid := range jwtKeypairs {
				keypair = jwtKeypairs[kid] // picks the last keypair...
			}
		}

		if keypair != nil {
			return &keypair.PublicKey, nil
		}

		return nil, errors.New("failed to resolve a valid JWT verification key")
	})
}

func init() {
	apiTokensInitCmd.Flags().StringVar(&common.ApplicationID, "application", "", "application id")
	apiTokensInitCmd.Flags().StringVar(&common.OrganizationID, "organization", "", "organization id")

	apiTokensInitCmd.Flags().BoolVar(&offlineAccess, "offline-access", false, "offline access")
	apiTokensInitCmd.Flags().BoolVar(&refreshToken, "refresh-token", false, "refresh token")
	apiTokensInitCmd.Flags().BoolVarP(&optional, "optional", "", false, "List all the optional flags")
}
