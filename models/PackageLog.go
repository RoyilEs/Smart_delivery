package models

import "gorm.io/gorm"

type Package struct {
	gorm.Model
	LogisticsId string `json:"logisticsId" gorm:"type:varchar(128);"`
	Action      string `json:"action" gorm:"type:varchar(128);"`
	Operator    string `json:"operator" gorm:"type:varchar(128);"`
	CreatedAt   string `json:"created_at" gorm:"type:varchar(128);"`
	Detail      string `json:"detail" gorm:"type:varchar(128);"`
	Status      string `json:"status" gorm:"column:status;type:varchar(50);comment:状态"`
	OperatorId  string `json:"operator_id" gorm:"column:operator_id;type:varchar(64);comment:操作人ID"`
	IpAddress   string `json:"ip_address" gorm:"column:ip_address;type:varchar(50);comment:操作IP"`
}

func (Package) TableName() string {
	return "packages"
}

// 操作类型常量
const (
	PackageActionCreate         = "create"          // 创建订单
	PackageActionUpdate         = "update"          // 更新订单
	PackageActionStore          = "store"           // 存入柜子
	PackageActionPickup         = "pickup"          // 取件
	PackageActionExpire         = "expire"          // 过期
	PackageActionAutoOutbound   = "auto_outbound"   // 自动出库
	PackageActionManualOutbound = "manual_outbound" // 手动出库
	PackageActionDelete         = "delete"          // 删除订单
	PackageActionNotify         = "notify"          // 发送通知
)
