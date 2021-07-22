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
	"github.com/gin-gonic/gin"
	"github.com/vzau/sso/middleware"
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
