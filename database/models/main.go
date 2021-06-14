package models

import (
	"fmt"
	"time"

	"github.com/dhawton/log4g"
	"gitlab.com/kzdv/sso/database/seed"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var MaxAttempts = 10
var DelayBetweenAttempts = time.Minute * 1
var attempt = 1

var log = log4g.Category("db")

func Connect(user string, pass string, hostname string, port string, database string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, pass, hostname, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error("Error connecting to database: " + err.Error())
		if attempt < MaxAttempts {
			log.Info(fmt.Sprintf("Attempt %d/%d Failed. Waiting %s before trying again...", attempt, MaxAttempts, DelayBetweenAttempts.String()))
			time.Sleep(DelayBetweenAttempts)
			attempt += 1
			Connect(user, pass, hostname, port, database)
			return
		}
		panic("Max attempts occured. Aborting startup.")
	}

	db.AutoMigrate(&OAuthClient{}, &OAuthLogin{}, &Rating{}, &User{})
	seed.CheckSeeds()

	DB = db
}
