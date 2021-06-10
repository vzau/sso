package main

import (
	"net/http"

	v1 "github.com/ZDV-Web-Team/sso/controllers/v1"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "PONG"})
	})

	engine.GET("/error", func(c *gin.Context) {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"message": GetPolicy().Sanitize(c.Query("message")),
		})
	})

	v1Router := engine.Group("/v1")
	{
		v1Router.GET("/authorize", v1.GetAuthorize)
		v1Router.GET("/return", v1.GetAuthorize)
		v1Router.GET("/token", v1.GetAuthorize)
	}
}
