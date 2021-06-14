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
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gitlab.com/kzdv/sso/database/models"
)

type Result struct {
	cid int
	err error
}

type UserResponse struct {
	CID int `json:"cid"`
}

func GetCallback(c *gin.Context) {
	token := c.Param("token")
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
		key, _ := jwk.Parse([]byte(os.Getenv("VATUSA_ULS_JWK")))
		_, err := jwt.Parse([]byte(token), jwt.WithKeySet(key), jwt.WithValidate(true))
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

	user := &models.User{}
	if err = models.DB.Where("cid = ?", userResult.cid).First(&user).Error; err != nil {
		handleError(c, "You are not part of our roster, so you are unable to login.")
		return
	}

	login.CID = user.CID
	login.Code, _ = gonanoid.New(32)

	c.Redirect(302, fmt.Sprintf("%s?code=%s&state=%s", login.RedirectURI, login.Code, login.State))
}
