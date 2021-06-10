package v1

import (
	"net/http"

	"github.com/ZDV-Web-Team/sso/lib/database/models"
	"github.com/gin-gonic/gin"
)

type AuthorizeRequest struct {
	ClientId            string `form:"client_id"`
	RedirectURI         string `form:"redirect_uri"`
	ResponseType        string `form:"response_type"`
	Scope               string `form:"scope" validation:"-"`
	CodeChallengeMethod string `form:"code_challenge_method"`
	CodeChallenge       string `form:"code_challenge"`
}

func GetAuthorize(c *gin.Context) {
	client := models.OAuthClient{}
	models.DB.Where("ClientId").First(&client)
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
