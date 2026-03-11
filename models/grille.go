package models

import (
	"Smart_delivery_locker/models/ctype"
	"gorm.io/gorm"
)

type Grille struct {
	gorm.Model
	GrilleId    string     `json:"grille_id" gorm:"type:varchar(128);"` // xx_xxx 拼接
	X           float64    `json:"x" gorm:"type:float;"`
	Y           float64    `json:"y" gorm:"type:float;"`
	Z           float64    `json:"z" gorm:"type:float;"`
	LogisticsId string     `json:"logisticsId" gorm:"type:varchar(128);"` // 订单号 所放入的订单
	Size        ctype.Size `json:"size" gorm:"type:int(10);"`
}
