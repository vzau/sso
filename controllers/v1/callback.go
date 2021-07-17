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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/vchicago/sso/database/models"
	dbTypes "github.com/vchicago/types/database"
)

type Result struct {
	cid int
	err error
}

type UserResponse struct {
	CID int `json:"cid"`
}

func GetCallback(c *gin.Context) {
	if _, cancel := c.GetQuery("cancel"); cancel {
		handleError(c, "Authentication cancelled.")
		return
	}

	token, exists := c.GetQuery("token")
	if !exists || len(token) < 1 {
		handleError(c, "Invalid response received from Authenticator or Authentication cancelled.")
		return
	}

	cookie, err := c.Cookie("sso_token")
	if err != nil {
		log4g.Category("controllers/callback").Error("Could not parse sso_token cookie, expired? " + err.Error())
		handleError(c, "Could not parse session cookie. Is it expired?")
		return
	}

	login := dbTypes.OAuthLogin{}
	if err = models.DB.Where("token = ? AND created_at < ?", cookie, time.Now().Add(time.Minute*5)).First(&login).Error; err != nil {
		log4g.Category("controllers/callback").Error("Token used that isn't in db, duplicate request? " + cookie)
		handleError(c, "Token is invalid.")
		return
	}

	if login.UserAgent != c.Request.UserAgent() {
		handleError(c, "Token is not valid.")
		go models.DB.Delete(login)
		return
	}

	result := make(chan Result)
	go func() {
		key, _ := jwk.ParseKey([]byte(os.Getenv("ULS_JWK")))
		_, err := jwt.Parse([]byte(token), jwt.WithVerify(jwa.SignatureAlgorithm(key.Algorithm()), key), jwt.WithValidate(true))
		if err != nil {
			log4g.Category("controllers/callback").Error("Error getting token information from VATUSA: " + err.Error())
			result <- Result{cid: 0, err: err}
			return
		}

		userdata := UserResponse{}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://login.vatusa.net/uls/v2/info?token=%s", token), bytes.NewBuffer(nil))
		req.Header.Add("Accept", "application/json")

		client := &http.Client{}
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			for key, val := range via[0].Header {
				req.Header[key] = val
			}

			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			log4g.Category("controllers/callback").Error("Error getting user information from VATUSA: " + err.Error())
			result <- Result{cid: 0, err: err}
			return
		}
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)

		if err = json.Unmarshal(data, &userdata); err != nil {
			log4g.Category("controllers/callback").Error("Error unmarshalling user data from VATUSA: " + string(data) + "--" + err.Error())
			result <- Result{cid: 0, err: err}
			return
		}

		result <- Result{cid: userdata.CID, err: err}
	}()

	userResult := <-result

	if userResult.err != nil {
		handleError(c, "Internal Error while getting user data from VATUSA Connect")
		return
	}

	user := &dbTypes.User{}
	if err = models.DB.First(&user, userResult.cid).Error; err != nil {
		handleError(c, "You are not part of our roster, so you are unable to login.")
		return
	}

	login.CID = user.CID
	login.Code, _ = gonanoid.New(32)
	models.DB.Save(&login)

	c.Redirect(302, fmt.Sprintf("%s?code=%s&state=%s", login.RedirectURI, login.Code, login.State))
}
