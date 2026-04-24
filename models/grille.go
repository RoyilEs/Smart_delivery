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
	Size        ctype.Size `json:"size_type" gorm:"type:int(10);"`

	CabinetId   string `json:"cabinet_id" gorm:"type:varchar(128);"`
	CabinetCode string `json:"cabinet_code" gorm:"type:varchar(128);"`

	MatrixRow    int `json:"matrix_row" gorm:"type:int(10);"`
	MatrixColumn int `json:"matrix_column" gorm:"type:int(10);"`
	Layer        int `json:"layer" gorm:"type:int(10);"`

	Remark string `json:"remark" gorm:"type:varchar(128);"`
	Status string `json:"status" gorm:"type:varchar(128);"`
}
