package seed

import (
	"errors"

	"github.com/dhawton/log4g"
	"github.com/vchicago/sso/database/models"
	dbTypes "github.com/vchicago/types/database"
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
