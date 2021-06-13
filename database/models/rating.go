package models

type Rating struct {
	ID    int    `json:"id"`
	Long  string `json:"long" gorm:"type:varchar(25)"`
	Short string `json:"short" gorm:"type:varchar(3)"`
}
