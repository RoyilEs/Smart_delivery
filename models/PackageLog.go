package models

import "gorm.io/gorm"

type Package struct {
	gorm.Model
	LogisticsId string `json:"logisticsId" gorm:"type:varchar(128);"`
	Action      string `json:"action" gorm:"type:varchar(128);"`
	Operator    string `json:"operator" gorm:"type:varchar(128);"`
	CreatedAt   string `json:"created_at" gorm:"type:varchar(128);"`
	Detail      string `json:"detail" gorm:"type:varchar(128);"`
}
