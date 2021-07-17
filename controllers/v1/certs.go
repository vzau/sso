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
	"encoding/json"
	"net/http"
	"os"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
)

type certReturn struct {
	message string
	keys    jwk.Set
}

func GetCerts(c *gin.Context) {
	jkeyset := os.Getenv("SSO_JWKS")
	keyset, err := jwk.Parse([]byte(jkeyset))
	if err != nil {
		log4g.Category("controllers/certs").Error("Error parsing JWKs: " + err.Error())
		log4g.Category("controllers/certs").Debug(jkeyset)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse JWKs"})
		return
	}

	pub, err := jwk.PublicSetOf(keyset)
	if err != nil {
		log4g.Category("controllers/certs").Error("Error generating public keyset: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate public keyset"})
		return
	}
	jpub, _ := json.Marshal(pub)
	out := map[string]interface{}{}
	json.Unmarshal([]byte(jpub), &out)
	out["message"] = "OK"
	c.JSON(http.StatusOK, out)
}
