/*
   ZAU Single Sign-On
   Copyright (C) 2021  Daniel A. Hawton <daniel@hawton.org>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package v1

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/vzau/sso/database/models"
	"github.com/vzau/sso/utils"
	dbTypes "github.com/vzau/types/database"
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

	client := dbTypes.OAuthClient{}
	if err := models.DB.Where("client_id = ?", req.ClientId).First(&client).Error; err != nil {
		handleError(c, "Invalid Client ID Received.")
		return
	}

	if ok, _ := client.ValidURI(req.RedirectURI); !ok {
		log4g.Category("controllers/authorize").Error("Unauthorized redirect uri received from client " + client.ClientId + ", " + req.RedirectURI)
		handleError(c, "The Return URI was not authorized.")
		return
	}

	if req.ResponseType != "code" {
		log4g.Category("controllers/authorize").Error("Invalid response type received from client " + client.ClientId + ", " + req.ResponseType)
		handleError(c, "Unsupported response type received.")
		return
	}

	if req.CodeChallengeMethod != "" && req.CodeChallengeMethod != "S256" {
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

	login := dbTypes.OAuthLogin{
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

	redirect_url := fmt.Sprintf("https://login.vatusa.net/uls/v2/login?fac=%s&url=%s&rfc7519_compliance", utils.Getenv("ULS_FACILITY_ID", "ZAU"), utils.Getenv("ULS_REDIRECT_ID", "1"))

	/*
		redirect_uri := url.QueryEscape(os.Getenv("VATSIM_REDIRECT_URI"))
		vatsim_url := fmt.Sprintf("https://auth.vatsim.net/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&response_type=code", os.Getenv("VATSIM_OAUTH_CLIENT_ID"), redirect_uri, url.QueryEscape("full_name email vatsim_details country"))
	*/
	c.Redirect(http.StatusTemporaryRedirect, redirect_url)
}
