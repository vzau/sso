package models

type OAuthLogin struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	Token     string      `json:"token" gorm:"type:varchar(128)"`
	UserAgent string      `json:"ua" gorm:"type:varchar(255)"`
	ReturnURL string      `json:"url" gorm:"type:varchar(255)"`
	ClientID  uint        `json:"-"`
	Client    OAuthClient `json:"-"`
}
