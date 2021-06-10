package models

import "gorm.io/datatypes"

type OAuthClient struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	ClientId     string            `json:"client_id" gorm:"type:varchar(128)"`
	ClientSecret string            `json:"-" gorm:"type:varchar(255)"`
	ReturnURIs   datatypes.JSONMap `json:"return_uris"`
}
