package v1

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kzdv/sso/database/models"
)

type TokenRequest struct {
	GrantType    string `form:"grant_type"`
	ClientId     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	CodeVerifier string `form:"code_verifier"`
}

type TokenResponse struct {
	AccessToken         string   `json:"access_token"`
	ExpiresIn           int      `json:"expires_in"`
	Scope               []string `json:"scope"`
	TokenType           string   `json:"token_type"`
	CodeChallenge       string   `json:"code_challenge"`
	CodeChallengeMethod string   `json:"code_challenge_method"`
}

func PostToken(c *gin.Context) {
	treq := TokenRequest{}
	if err := c.ShouldBind(&treq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	if treq.GrantType != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	login := models.OAuthLogin{}
	if err := models.DB.Where("code = ?", treq.Code).First(&login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	defer models.DB.Delete(&login)

	if treq.ClientId != login.Client.ClientId || treq.ClientSecret != login.Client.ClientSecret {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client"})
		return
	}

	hash := sha256.Sum256([]byte(treq.CodeVerifier))
	if login.CodeChallenge != base64.RawURLEncoding.EncodeToString(hash[:]) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}
}
