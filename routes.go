package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "gitlab.com/kzdv/sso/controllers/v1"
)

func SetupRoutes(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "PONG"})
	})

	v1Router := engine.Group("/v1")
	{
		v1Router.GET("/authorize", v1.GetAuthorize)
		v1Router.GET("/callback", v1.GetCallback)
		v1Router.GET("/certs", v1.GetCerts)
		v1Router.GET("/token", v1.GetAuthorize)
	}
}
