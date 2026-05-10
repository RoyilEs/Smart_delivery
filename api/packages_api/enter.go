package packages_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/res"
	"github.com/gin-gonic/gin"
)

type PackageApi struct{}

type PackageListQuery struct {
	Keyword string `json:"keyword" form:"keyword"` // 按物流单号、收件人、手机号搜索
	Status  string `json:"status" form:"status"`
}

func (PackageApi) PackageListView(c *gin.Context) {
	var query PackageListQuery
	err := c.BindQuery(&query)
	if err != nil {
		res.ResultFailWithError(err, &query, c)
		return
	}
	var cr []models.Item
	db := global.DB.Model(&models.Item{})

	// 关键词搜索
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("logistics_id LIKE ? OR receiver_name LIKE ? OR receiver_phone LIKE ?",
			keyword, keyword, keyword)
	}

	// 状态筛选
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	err = db.Debug().Find(&cr).Error
	if err != nil {
		res.ResultFailWithMsg("包裹列表获取失败", c)
		return
	}
	res.ResultOkWithList(cr, int64(len(cr)), c)
}
