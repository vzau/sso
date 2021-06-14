package seed

import (
	"encoding/json"

	"github.com/dhawton/log4g"
	"gitlab.com/kzdv/sso/database/models"
	"gorm.io/gorm/clause"
)

type RatingInfo struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func SeedRating() {
	var ratings = `
		[
			{"short":"OBS","long":"Observer"},
			{"short":"S1","long":"Student 1"},
			{"short":"S2","long":"Student 2"},
			{"short":"S3","long":"Student 3"},
			{"short":"C1","long":"Controller"},
			{"short":"C2","long":"Controller 2"},
			{"short":"C3","long":"Controller 3"},
			{"short":"I1","long":"Instructor"},
			{"short":"I2","long":"Instructor 2"},
			{"short":"I3","long":"Senior Instructor"},
			{"short":"SUP","long":"Supervisor"},
			{"short":"ADM","long":"Administrator"},
		]
	`

	var ratingsDecoded []RatingInfo
	err := json.Unmarshal([]byte(ratings), &ratingsDecoded)
	if err != nil {
		log4g.Category("SeedRating").Error("Could not decode ratings for seeding: " + err.Error())
	}

	for k, v := range ratingsDecoded {
		rating := &models.Rating{
			ID:    k,
			Long:  v.Long,
			Short: v.Short,
		}

		models.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"short", "long"}),
		}).Create(&rating)
	}
}
