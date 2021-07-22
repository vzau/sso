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
	"fmt"
	"os"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/vzau/sso/database/models"
	"github.com/vzau/sso/database/seed"
	"github.com/vzau/sso/utils"
	dbTypes "github.com/vzau/types/database"
)

var log = log4g.Category("main")

func main() {
	log4g.SetLogLevel(log4g.DEBUG)

	intro := figure.NewFigure("ZDV SSO", "", false).Slicify()
	for i := 0; i < len(intro); i++ {
		log.Info(intro[i])
	}

	log.Info("Starting ZDV SSO")
	log.Info("Checking for .env, loading if exists")
	if _, err := os.Stat(".env"); err == nil {
		log.Info("Found, loading")
		err := godotenv.Load()
		if err != nil {
			log.Error("Error loading .env file: " + err.Error())
		}
	}

	appenv := utils.Getenv("APP_ENV", "dev")
	log.Debug(fmt.Sprintf("APPENV=%s", appenv))

	if appenv == "production" {
		log4g.SetLogLevel(log4g.INFO)
		log.Info("Setting gin to Release Mode")
		gin.SetMode(gin.ReleaseMode)
	} else {
		log4g.SetLogLevel(log4g.DEBUG)
	}

	log.Info("Connecting to database and handling migrations")
	models.Connect(utils.Getenv("DB_USERNAME", "root"), utils.Getenv("DB_PASSWORD", "secret"), utils.Getenv("DB_HOSTNAME", "localhost"), utils.Getenv("DB_PORT", "3306"), utils.Getenv("DB_DATABASE", "zdv"))
	seed.CheckSeeds()

	log.Info("Configuring Gin Server")
	server := NewServer(appenv)

	log.Info("Configuring scheduled jobs")
	jobs := cron.New()
	jobs.AddFunc("@every 1m", func() {
		if err := models.DB.Where("created_at >= ?", time.Now().Add(time.Minute*30).Unix()).Delete(&dbTypes.OAuthLogin{}).Error; err != nil {
			log4g.Category("job/cleanup").Error(fmt.Sprintf("Error cleaning up old codes: %s", err.Error()))
		}
	})
	jobs.Start()

	log.Info("Done with setup, starting web server...")
	server.engine.Run(fmt.Sprintf(":%s", utils.Getenv("PORT", "3000")))
}
