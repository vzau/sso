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

package seed

import (
	"errors"

	"github.com/dhawton/log4g"
	"github.com/vzau/sso/database/models"
	dbTypes "github.com/vzau/types/database"
	"gorm.io/gorm"
)

var log = log4g.Category("seed")

func CheckSeeds() {
	// Check if Ratings should be seeded
	log.Debug("Checking ratings")
	var r = dbTypes.Rating{}
	if err := models.DB.Where("ID = ?", 1).First(&r).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug("Check failed for Record Not Found, seeding Ratings")
			SeedRating()
		}
	}
}
