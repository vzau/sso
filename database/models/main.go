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

package models

import (
	"fmt"
	"time"

	"github.com/dhawton/log4g"
	dbTypes "github.com/vzau/types/database"
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

	db.AutoMigrate(&dbTypes.OAuthClient{}, &dbTypes.OAuthLogin{}, &dbTypes.Rating{}, &dbTypes.Role{}, &dbTypes.User{})

	DB = db
}
