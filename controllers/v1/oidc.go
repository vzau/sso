package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OIDCConfig struct {
	Issuer                                     string   `json:"issuer"`
	AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
	TokenEndpoint                              string   `json:"token_endpoint"`
	UserinfoEndpoint                           string   `json:"userinfo_endpoint"`
	JwksUri                                    string   `json:"jwks_uri"`
	GrantTypesSupported                        []string `json:"grant_types_supported"`
	ResponseTypesSupported                     []string `json:"response_types_supported"`
	IdTokenSigningAlgValuesSupported           []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported"`
	CodeChallengeMethodsSupported              []string `json:"code_challenge_methods_supported"`
	RequestParameterSupported                  bool     `json:"request_parameter_supported"`
	RequestUriParameterSupported               bool     `json:"request_uri_parameter_supported"`
	ScopesSupported                            []string `json:"scopes_supported"`
	ClaimsSupported                            []string `json:"claims_supported"`
}

func GetOIDCConfig(c *gin.Context) {
	host := c.Request.Host
	config := OIDCConfig{
		Issuer:                 "https://" + host,
		AuthorizationEndpoint:  "https://" + host + "/oauth/authorize",
		TokenEndpoint:          "https://" + host + "/oauth/token",
		UserinfoEndpoint:       "https://" + host + "/v1/info",
		JwksUri:                "https://" + host + "/oauth/certs",
		GrantTypesSupported:    []string{"authorization_code"},
		ResponseTypesSupported: []string{"code"},
		IdTokenSigningAlgValuesSupported: []string{
			"PS384", "ES384", "RS384", "HS256", "HS512", "ES256", "RS256", "HS384", "ES512", "PS256", "PS512", "RS512",
		},
		TokenEndpointAuthSigningAlgValuesSupported: []string{
			"PS384", "ES384", "RS384", "HS256", "HS512", "ES256", "RS256", "HS384", "ES512", "PS256", "PS512", "RS512",
		},
		CodeChallengeMethodsSupported: []string{"none", "S256"},
		RequestParameterSupported:     true,
		RequestUriParameterSupported:  true,
		ScopesSupported:               []string{"openid", "profile", "email"},
		ClaimsSupported: []string{
			"sub",
			"name",
			"iss",
			"exp",
			"iat",
		},
	}

	c.JSON(http.StatusOK, config)
}
