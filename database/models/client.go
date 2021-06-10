package models

import (
	"encoding/json"
	"fmt"

	"github.com/dhawton/log4g"
)

type OAuthClient struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	ClientId     string `json:"client_id" gorm:"type:varchar(128)"`
	ClientSecret string `json:"-" gorm:"type:varchar(255)"`
	RedirectURIs string `json:"return_uris" gorm:"type:text"`
}

func (c *OAuthClient) ValidURI(uri string) bool {
	uris := []string{}
	err := json.Unmarshal([]byte(c.RedirectURIs), &uris)
	if err != nil {
		log4g.Category("model/client/ValidURI").Error(fmt.Sprintf("Error unmarshalling RedirectURIs: %s", err.Error()))
		return false
	}
	for _, v := range uris {
		if uri == v {
			return true
		}
	}

	return false
}
