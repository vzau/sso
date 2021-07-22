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
	"encoding/json"

	"github.com/dhawton/log4g"
	"github.com/vzau/sso/database/models"
	dbTypes "github.com/vzau/types/database"
	"gorm.io/gorm/clause"
)

type RatingInfo struct {
	ID    int    `json:"id"`
	Short string `json:"short"`
	Long  string `json:"long"`
}

func SeedRating() {
	var ratings = `
		[
			{"id":-1,"short":"INA","long":"Inactive"},
			{"id":0,"short":"SUS","long":"Suspended"},
			{"id":1,"short":"OBS","long":"Observer"},
			{"id":2,"short":"S1","long":"Student 1"},
			{"id":3,"short":"S2","long":"Student 2"},
			{"id":4,"short":"S3","long":"Student 3"},
			{"id":5,"short":"C1","long":"Controller"},
			{"id":6,"short":"C2","long":"Controller 2"},
			{"id":7,"short":"C3","long":"Controller 3"},
			{"id":8,"short":"I1","long":"Instructor"},
			{"id":9,"short":"I2","long":"Instructor 2"},
			{"id":10,"short":"I3","long":"Senior Instructor"},
			{"id":11,"short":"SUP","long":"Supervisor"},
			{"id":12,"short":"ADM","long":"Administrator"}
		]
	`

	var ratingsDecoded []RatingInfo
	err := json.Unmarshal([]byte(ratings), &ratingsDecoded)
	if err != nil {
		log4g.Category("SeedRating").Error("Could not decode ratings for seeding: " + err.Error())
	}

	for _, v := range ratingsDecoded {
		rating := &dbTypes.Rating{
			ID:    v.ID,
			Long:  v.Long,
			Short: v.Short,
		}

		models.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"short", "long"}),
		}).Create(&rating)
	}
}
