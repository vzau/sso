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
	"gitlab.com/kzdv/sso/database/models"
	"gitlab.com/kzdv/sso/database/seed"
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

	appenv := Getenv("APP_ENV", "dev")
	log.Debug(fmt.Sprintf("APPENV=%s", appenv))

	if appenv == "production" {
		log4g.SetLogLevel(log4g.INFO)
		log.Info("Setting gin to Release Mode")
		gin.SetMode(gin.ReleaseMode)
	} else {
		log4g.SetLogLevel(log4g.DEBUG)
	}

	log.Info("Connecting to database and handling migrations")
	models.Connect(Getenv("DB_USERNAME", "root"), Getenv("DB_PASSWORD", "secret"), Getenv("DB_HOSTNAME", "localhost"), Getenv("DB_PORT", "3306"), Getenv("DB_DATABASE", "zdv"))
	seed.CheckSeeds()

	log.Info("Configuring Gin Server")
	server := NewServer(appenv)

	log.Info("Configuring scheduled jobs")
	jobs := cron.New()
	jobs.AddFunc("@every 5m", func() {
		if err := models.DB.Where("created_at >= ?", time.Now().Add(time.Minute*30).Unix()).Delete(&models.OAuthLogin{}).Error; err != nil {
			log4g.Category("job/cleanup").Error(fmt.Sprintf("Error cleaning up old codes: %s", err.Error()))
		}
	})
	jobs.Start()

	log.Info("Done with setup, starting web server...")
	server.engine.Run(fmt.Sprintf(":%s", Getenv("PORT", "3000")))
}
