package packages_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype/action"
	"Smart_delivery_locker/models/res"
	"Smart_delivery_locker/utils"
	"Smart_delivery_locker/utils/jwts"
	"github.com/gin-gonic/gin"
	"strconv"
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

func (PackageApi) PackageUpdateVIew(c *gin.Context) {
	id := c.Param("id")
	var cr models.Item
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}
	var item models.Item
	err = global.DB.Model(&item).Where("id = ?", id).Updates(cr).Error
	if err != nil {
		res.ResultFailWithMsg("包裹更新失败", c)
		return
	}
	global.DB.Find(&item, id)

	cl, err := jwts.ParseToken(c.GetHeader("token"))
	if err != nil {
		res.ResultFailWithMsg("token解码错误", c)
	}
	utils.RecordPackageLog(item.LogisticsId, action.Updated.String(), models.PackageActionUpdate, "包裹更新成功", cl.Username, strconv.Itoa(int(cl.UserID)), c)

	res.ResultOkWithData(item, c)
}
