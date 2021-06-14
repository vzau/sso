/* Moved to vatsim in case we need to revert back to VATSIM Connect later */

package vatsim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gitlab.com/kzdv/sso/database/models"
)

type OAuthResponse struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

type UserData struct {
	CID int `json:"cid"`
}

type UserResponse struct {
	Data UserData `json:"data"`
}

type Result struct {
	cid int
	err error
}

func GetCallback(c *gin.Context) {
	token := c.Param("code")
	if len(token) < 1 {
		c.HTML(http.StatusInternalServerError, "error.tmpl", "Invalid response received from Authenticator or Authentication cancelled.")
		return
	}

	cookie, err := c.Cookie("sso_token")
	if err != nil {
		log4g.Category("controllers/callback").Error("Could not parse sso_token cookie, expired? " + err.Error())
		c.HTML(http.StatusInternalServerError, "error.tmpl", "Could not parse session cookie. Is it expired?")
		return
	}

	login := models.OAuthLogin{}
	if err = models.DB.Where("token = ? AND created_at < ?", cookie, time.Now().Add(time.Minute*5)).First(&login).Error; err != nil {
		log4g.Category("controllers/callback").Error("Token used that isn't in db, duplicate request? " + cookie)
		c.HTML(http.StatusInternalServerError, "error.tmpl", "Token is invalid.")
		return
	}

	if login.UserAgent != c.Request.UserAgent() {
		handleError(c, "Token is not valid.")
		go models.DB.Delete(login)
		return
	}

	result := make(chan Result)
	go func() {
		resp, err := http.Post(
			fmt.Sprintf(
				"https://auth.vatsim.net/oauth/token?grant_type=authorization_code&client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
				os.Getenv("VATSIM_OAUTH_CLIENT_ID"),
				os.Getenv("VATSIM_OAUTH_CLIENT_SECRET"),
				url.QueryEscape(os.Getenv("VATSIM_REDIRECT_URI")),
				token,
			), "application/json", bytes.NewBuffer(nil))
		if err != nil {
			log4g.Category("controllers/callback").Error("Error getting token information from VATSIM: " + err.Error())
			result <- Result{cid: 0, err: err}
			return
		}

		oauthresponse := OAuthResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &oauthresponse); err != nil {
			log4g.Category("controllers/callback").Error("Error parsing JSON object from VATSIM: " + string(body) + " -- " + err.Error())
			result <- Result{cid: 0, err: err}
			return
		}

		userdata := UserResponse{}
		req, err := http.NewRequest("GET", "https://auth.vatsim.net/api/user", bytes.NewBuffer(nil))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauthresponse.AccessToken))
		req.Header.Add("Accept", "application/json")

		client := &http.Client{}
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			for key, val := range via[0].Header {
				req.Header[key] = val
			}

			return err
		}
		resp, err = client.Do(req)
		if err != nil {
			log4g.Category("controllers/callback").Error("Error getting user information from VATSIM: " + err.Error())
			result <- Result{cid: 0, err: err}
			return
		}
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)

		if err = json.Unmarshal(data, &userdata); err != nil {
			log4g.Category("controllers/callback").Error("Error unmarshalling user data from VATSIM: " + string(data) + "--" + err.Error())
			result <- Result{cid: 0, err: err}
			return
		}

		result <- Result{cid: userdata.Data.CID, err: err}
	}()

	userResult := <-result

	if userResult.err != nil {
		handleError(c, "Internal Error while getting user data from VATSIM Connect")
		return
	}

	user := &models.User{}
	if err = models.DB.Where("cid = ?", userResult.cid).First(&user).Error; err != nil {
		handleError(c, "You are not part of our roster, so you are unable to login.")
		return
	}

	login.CID = user.CID
	login.Code, _ = gonanoid.New(32)

	c.Redirect(302, fmt.Sprintf("%s?code=%s&state=%s", login.RedirectURI, login.Code, login.State))
}
