package models

import "gorm.io/gorm"

// Item 包裹信息表
type Item struct {
	gorm.Model
	ReceiverName    string  `json:"receiverName" gorm:"type:varchar(40);not null;"`
	ReceiverPhone   string  `json:"receiverPhone" gorm:"type:varchar(20);"`
	ReceiverEmail   string  `json:"receiverEmail" gorm:"type:varchar(20);"`
	ReceiverCity    string  `json:"receiverCity" gorm:"type:varchar(20);"`
	ReceiverArea    string  `json:"receiverArea" gorm:"type:varchar(20);"`
	ReceiverAddress string  `json:"receiverAddress" gorm:"type:varchar(256);"`
	SenderName      string  `json:"senderName" gorm:"type:varchar(40);not null;"`
	SenderPhone     string  `json:"senderPhone" gorm:"type:varchar(20);"`
	SenderEmail     string  `json:"senderEmail" gorm:"type:varchar(20);"`
	SenderCity      string  `json:"senderCity" gorm:"type:varchar(20);"`
	SenderArea      string  `json:"senderArea" gorm:"type:varchar(20);"`
	SenderAddress   string  `json:"senderAddress" gorm:"type:varchar(256);"`
	ItemName        string  `json:"itemName" gorm:"type:varchar(128);"`
	ItemNum         int     `json:"itemNum" gorm:"type:int(10);"`          // 商品数量
	ItemWeight      float64 `json:"itemWeight" gorm:"type:float;"`         // 商品重量
	PackageNums     int     `json:"packageNums" gorm:"type:int(10);"`      // 包裹数量
	LogisticsId     string  `json:"logisticsId" gorm:"type:varchar(128);"` // 订单号
}
