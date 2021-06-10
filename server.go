package main

import (
	"github.com/ZDV-Web-Team/sso/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(appenv string) *Server {
	server := Server{}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger)
	server.engine = engine
	engine.LoadHTMLGlob("templates/*")

	SetupRoutes(engine)

	return &server
}
