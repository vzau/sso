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

package jwt

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/vzau/sso/database/models"
	dbTypes "github.com/vzau/types/database"
	"gorm.io/gorm/clause"
)

var (
	ErrNoToken = errors.New("No Token Specified")
)

var requireAuth bool
var log = log4g.Category("middleware/jwt")

func Auth(c *gin.Context) {
	requireAuth = true

	const BEARER_SCHEMA = "Bearer"
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < len(BEARER_SCHEMA) {
		HandleRet(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokenString := authHeader[len(BEARER_SCHEMA):]
	keyset, err := jwk.Parse([]byte(os.Getenv("SSO_JWKS")))
	if err != nil {
		HandleRet(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(keyset), jwt.WithValidate(true))
	if err != nil {
		log.Warning("Bad token passed: %s // %s", err.Error(), tokenString)
		HandleRet(c, http.StatusForbidden, "Forbidden")
		return
	}

	cid, err := strconv.ParseUint(token.Subject(), 10, 32)
	if err != nil {
		log.Warning("Cannot convert subject to int, bad! %s // %s // %s", err.Error(), token.Subject(), tokenString)
		HandleRet(c, http.StatusForbidden, "Forbidden")
		return
	}

	user := &dbTypes.User{}
	if err = models.DB.Where(&dbTypes.User{CID: uint(cid)}).Preload(clause.Associations).First(&user).Error; err != nil {
		log.Warning("No user found for %d: %s", cid, err.Error())
		HandleRet(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Set("x-user", user)
	c.Next()
}

func HandleRet(c *gin.Context, ret int, msg string) {
	if !requireAuth {
		c.Set("x-user", nil)
		c.Next()
	} else {
		c.JSON(ret, gin.H{"message": msg})
		c.Abort()
	}
}
