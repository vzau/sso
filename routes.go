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

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/vzau/sso/controllers/v1"
	jwtMiddleware "github.com/vzau/sso/middleware/jwt"
)

func SetupRoutes(engine *gin.Engine) {
	engine.StaticFile("/", "./templates/docs.html")
	engine.StaticFile("/openapi.json", "./templates/openapi.json")

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "PONG"})
	})

	OAuthRouter := engine.Group("/oauth")
	{
		OAuthRouter.GET("/authorize", v1.GetAuthorize)
		OAuthRouter.GET("/callback", v1.GetCallback)
		OAuthRouter.GET("/certs", v1.GetCerts)
		OAuthRouter.POST("/token", v1.PostToken)
	}

	v1Router := engine.Group("/v1")
	{
		v1Router.GET("/info", jwtMiddleware.Auth, v1.GetInfo)
	}
}
