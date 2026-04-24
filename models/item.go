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
	X               float64 `json:"x" gorm:"type:float;"`
	Y               float64 `json:"y" gorm:"type:float;"`
	Z               float64 `json:"z" gorm:"type:float;"`
	GrilleId        string  `json:"grille_id" gorm:"type:varchar(128);"`

	PickupCOde string `json:"pickup_code" gorm:"type:varchar(128);"` // 取件码
	Status     string `json:"status" gorm:"type:varchar(128);"`      // 包裹状态

	CabinetId   string `json:"cabinet_id" gorm:"type:varchar(128);"`   // 柜体ID
	CabinetCode string `json:"cabinet_code" gorm:"type:varchar(128);"` // 柜体名称

	GrilleStatus  string `json:"grille_status" gorm:"type:varchar(128);"` // 当前格口状态
	InboundAt     string `json:"inbound_at" gorm:"type:varchar(128);"`    // 入柜时间
	OutboundAt    string `json:"outbound_at" gorm:"type:varchar(128);"`   // 出柜时间
	Remark        string `json:"remark" gorm:"type:varchar(128);"`        // 备注
	ReceiverToken string `json:"receiver_token" gorm:"type:varchar(128);"`
}
