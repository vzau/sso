package v1

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
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
	AccessToken         string `json:"access_token"`
	ExpiresIn           int    `json:"expires_in"`
	TokenType           string `json:"token_type"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
}

func PostToken(c *gin.Context) {
	params, _ := json.Marshal(c.Params)
	log4g.Category("controllers/authorize").Debug(string(params))
	form, _ := json.Marshal(c.Request.PostForm)
	log4g.Category("controllers/authorize").Debug(string(form))

	treq := TokenRequest{}
	if err := c.ShouldBind(&treq); err != nil {
		log4g.Category("controllers/token").Error("Invalid request, missing field(s)")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	if treq.GrantType != "authorization_code" {
		log4g.Category("controllers/token").Error("Grant type is not authorization code")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	login := models.OAuthLogin{}
	if err := models.DB.Where("code = ?", treq.Code).First(&login).Error; err != nil {
		log4g.Category("controllers/token").Error(fmt.Sprintf("Code %s not found", treq.Code))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	//defer models.DB.Delete(&login)

	if treq.ClientId != login.Client.ClientId || treq.ClientSecret != login.Client.ClientSecret {
		log4g.Category("controllers/token").Error(fmt.Sprintf("Invalid client: %s %s does not match %s %s", treq.ClientId, treq.ClientSecret, login.Client.ClientId, login.Client.ClientSecret))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client"})
		return
	}

	hash := sha256.Sum256([]byte(treq.CodeVerifier))
	if login.CodeChallenge != base64.RawURLEncoding.EncodeToString(hash[:]) {
		log4g.Category("controllers/token").Error(fmt.Sprintf("Code Challenge failed"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	keyset, err := jwk.Parse([]byte(os.Getenv("ZDK_JWKS")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log4g.Category("controller/token").Error("Could not parse JWKs: " + err.Error())
		return
	}

	key, _ := keyset.LookupKeyID(os.Getenv("ZDV_CURRENT_KEY"))
	token := jwt.New()
	token.Set(jwt.IssuerKey, "sso.kzdv.io")
	token.Set(jwt.AudienceKey, login.Client.Name)
	token.Set(jwt.SubjectKey, fmt.Sprint(login.CID))
	token.Set(jwt.IssuedAtKey, time.Now())
	token.Set(jwt.ExpirationKey, (time.Hour*24*7)/time.Second)
	signed, err := jwt.Sign(token, jwa.SignatureAlgorithm(key.KeyType()), key)

	ret := TokenResponse{
		AccessToken:         string(signed),
		ExpiresIn:           int(time.Now().Add(time.Hour * 24 * 7).Unix()),
		TokenType:           "Bearer",
		CodeChallenge:       login.CodeChallenge,
		CodeChallengeMethod: login.CodeChallengeMethod,
	}

	c.JSON(http.StatusOK, ret)
}
