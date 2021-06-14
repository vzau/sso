package v1

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gitlab.com/kzdv/sso/database/models"
)

type AuthorizeRequest struct {
	ClientId            string `form:"client_id"`
	RedirectURI         string `form:"redirect_uri"`
	ResponseType        string `form:"response_type"`
	Scope               string `form:"scope" validation:"-"`
	CodeChallengeMethod string `form:"code_challenge_method"`
	CodeChallenge       string `form:"code_challenge"`
	State               string `form:"state"`
}

func GetAuthorize(c *gin.Context) {
	req := AuthorizeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		handleError(c, "Invalid OAuth2 Request.")
		return
	}

	client := models.OAuthClient{}
	if err := models.DB.Where("client_id = ?", req.ClientId).First(&client).Error; err != nil {
		handleError(c, "Invalid Client ID Received.")
		return
	}

	if !client.ValidURI(req.RedirectURI) {
		log4g.Category("controllers/authorize").Error("Unauthorized redirect uri received from client " + client.ClientId + ", " + req.RedirectURI)
		handleError(c, "The Return URI was not authorized.")
		return
	}

	if req.ResponseType != "token" {
		log4g.Category("controllers/authorize").Error("Invalid response type received from client " + client.ClientId + ", " + req.ResponseType)
		handleError(c, "Unsupported response type received.")
		return
	}

	if req.CodeChallengeMethod != "S256" {
		log4g.Category("controllers/authorize").Error("Invalid code challenge method received from client " + client.ClientId + ", " + req.CodeChallengeMethod)
		handleError(c, "Unsupported Code Challenge Method defined.")
		return
	}

	token, err := gonanoid.New(32)
	if err != nil {
		log4g.Category("controllers/authorize").Error("Error generating new token " + err.Error())
		handleError(c, "Failed to generate new token.")
		return
	}

	login := models.OAuthLogin{
		Token:               token,
		UserAgent:           c.Request.UserAgent(),
		RedirectURI:         req.RedirectURI,
		Client:              client,
		ClientID:            client.ID,
		State:               req.State,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		Scope:               req.Scope,
	}

	if err = models.DB.Create(&login).Error; err != nil {
		log4g.Category("controllers/authorize").Error("Failed to store token " + err.Error())
		handleError(c, "Failed to create token")
		return
	}

	/*	Leave this here, hopefully ULSv3 lets us send the return url instead of an ID #

		scheme := "http"
		if c.Request.TLS != nil && c.Request.TLS.HandshakeComplete {
			scheme = "https"
		}
		returnUri := url.QueryEscape(fmt.Sprintf("%s://%s/v1/return", scheme, c.Request.Host)) */

	host, _, _ := net.SplitHostPort(c.Request.Host)
	if host == "" {
		host = c.Request.Host
	}
	log4g.Category("test").Debug(host)
	c.SetCookie("sso_token", login.Token, int(time.Minute)*5, "/", host, false, true)

	redirect_url := fmt.Sprintf("https://login.vatusa.net/uls/v2/login?fac=ZDV&url=%s&rfc7519_compliance", os.Getenv("VATUSA_ULS_REDIRECT_ID"))

	/*
		redirect_uri := url.QueryEscape(os.Getenv("VATSIM_REDIRECT_URI"))
		vatsim_url := fmt.Sprintf("https://auth.vatsim.net/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&response_type=code", os.Getenv("VATSIM_OAUTH_CLIENT_ID"), redirect_uri, url.QueryEscape("full_name email vatsim_details country"))
	*/
	c.Redirect(http.StatusTemporaryRedirect, redirect_url)
}
