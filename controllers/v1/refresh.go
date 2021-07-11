package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func GetRefresh(c *gin.Context) {
	/*	authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		keyset, _ := jwk.Parse([]byte(os.Getenv("SSO_JWKS")))

		t, err := jwt.Parse(parts[1], jwt.WithKeySet(keyset))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			return
		}

		key, _ := keyset.LookupKeyID(os.Getenv("SSO_CURRENT_KEY"))

		token := jwt.New()
		token.Set(jwt.IssuerKey, "sso.kzdv.io")
		token.Set(jwt.AudienceKey, login.Client.Name)
		token.Set(jwt.SubjectKey, string(login.CID))
		token.Set(jwt.IssuedAtKey, time.Now())
		token.Set(jwt.ExpirationKey, (time.Hour*24*7)/time.Second)
		signed, err := jwt.Sign(token, key.Algorithm(), key)

		ret := TokenResponse{
			AccessToken:         signed,
			ExpiresIn:           int(time.Now().Add(time.Hour * 24 * 7).Unix()),
			TokenType:           "Bearer",
			CodeChallenge:       login.CodeChallenge,
			CodeChallengeMethod: login.CodeChallengeMethod,
		} */

	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not Implemented"})
}
