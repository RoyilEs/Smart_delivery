package item_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype"
	"Smart_delivery_locker/models/res"
	CODE "Smart_delivery_locker/models/res/code"
	"Smart_delivery_locker/service/common"
	"Smart_delivery_locker/service/user_ser"
	"Smart_delivery_locker/utils/pwd"
	"github.com/gin-gonic/gin"
)

type ItemApi struct{}

// ItemUri 获取包裹信息
type ItemUri struct {
	Name string `uri:"name"`
}

type ItemResponse struct {
	models.Item
}

type ItemListRequest struct {
	models.PageInfo
	Permission int `json:"permission" form:"permission"`
}

func (ItemApi) ItemListView(c *gin.Context) {
	var cr ItemUri
	if err := c.ShouldBindUri(&cr); err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}

	var page ItemListRequest
	if err := c.ShouldBind(&page); err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}

	var userModel models.User
	err := global.DB.Where("username = ?", cr.Name).Find(&userModel).Error
	if err != nil {
		res.ResultFailWithMsg("用户不存在", c)
		return
	}

	var (
		items []ItemResponse
		count int64
	)
	// 普通用户
	if userModel.Permission == ctype.PermissionUser {
		list, _, _ := common.ComList(models.Item{SenderName: userModel.Username}, common.Option{
			PageInfo: page.PageInfo,
		})

		for _, item := range list {
			if item.SenderName == userModel.Username {
				items = append(items, ItemResponse{
					Item: item,
				})
				count += 1
			}
		}
	}
	// 快递员与管理员
	if userModel.Permission == ctype.PermissionCourier || userModel.Permission == ctype.PermissionAdmin {
		list, c, _ := common.ComList(models.Item{}, common.Option{
			PageInfo: page.PageInfo,
		})
		count = c
		for _, item := range list {
			items = append(items, ItemResponse{
				Item: item,
			})
		}
	}
	res.ResultOkWithList(items, count, c)
}

// ItemCreateRequest TODO 发货端口 是否仅为测试手动输入发货
type ItemCreateRequest struct {
	ReceiverName    string  `json:"receiverName" binding:"required" msg:"请输入收件人姓名"`
	ReceiverPhone   string  `json:"receiverPhone" binding:"required" msg:"请输入收件人手机号"`
	ReceiverEmail   string  `json:"receiverEmail" binding:"required" msg:"请输入收件人邮箱"`
	ReceiverCity    string  `json:"receiverCity" binding:"required" msg:"请输入收件人城市"`
	ReceiverArea    string  `json:"receiverArea" binding:"required" msg:"请输入收件人区县"`
	ReceiverAddress string  `json:"receiverAddress" binding:"required" msg:"请输入收件人地址"`
	SenderName      string  `json:"senderName" binding:"required" msg:"请输入发件人姓名"`
	SenderPhone     string  `json:"senderPhone" binding:"required" msg:"请输入发件人手机号"`
	SenderEmail     string  `json:"senderEmail" binding:"required" msg:"请输入发件人邮箱"`
	SenderCity      string  `json:"senderCity" binding:"required" msg:"请输入发件人城市"`
	SenderArea      string  `json:"senderArea" binding:"required" msg:"请输入发件人区县"`
	SenderAddress   string  `json:"senderAddress" binding:"required" msg:"请输入发件人地址"`
	ItemName        string  `json:"itemName" binding:"required" msg:"请输入商品名称"`
	ItemNum         int     `json:"itemNum" binding:"required" msg:"请输入商品数量"`
	ItemWeight      float64 `json:"itemWeight" binding:"required" msg:"请输入商品重量"`
	PackageNums     int     `json:"packageNums" binding:"required" msg:"请输入包裹数量"`
}

func (ItemApi) ItemCreateView(c *gin.Context) {
	var cr ItemCreateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}

	// 直接拼出对接码 为 双方手机号拼接
	logisticId := pwd.HashPassword(cr.ReceiverPhone + cr.SenderPhone)

	// 创建包裹
	item := models.Item{
		ItemName:        cr.ItemName,
		ItemNum:         cr.ItemNum,
		ItemWeight:      cr.ItemWeight,
		PackageNums:     cr.PackageNums,
		ReceiverName:    cr.ReceiverName,
		ReceiverPhone:   cr.ReceiverPhone,
		ReceiverEmail:   cr.ReceiverEmail,
		ReceiverCity:    cr.ReceiverCity,
		ReceiverArea:    cr.ReceiverArea,
		ReceiverAddress: cr.ReceiverAddress,
		SenderName:      cr.SenderName,
		SenderPhone:     cr.SenderPhone,
		SenderEmail:     cr.SenderEmail,
		SenderCity:      cr.SenderCity,
		SenderArea:      cr.SenderArea,
		SenderAddress:   cr.SenderAddress,
		LogisticsId:     logisticId,
	}

	// 检测这个邮寄用户是否在数据库中 不存在则建立
	var (
		receiverUser models.User
		senderUser   models.User
		count        int64
	)
	global.DB.Where("username = ?", cr.ReceiverName).Find(&receiverUser).Count(&count)
	if count == 0 && receiverUser.Phone != cr.ReceiverPhone {
		global.Log.Println("用户不存在直接建立")
		err := user_ser.UserService{}.CreateUser(cr.ReceiverName, "", ctype.PermissionUser, cr.ReceiverEmail, cr.ReceiverPhone)
		if err != nil {
			global.Log.Error(err)
			res.ResultFailWithMsg("用户创建失败!", c)
			return
		}
	}
	global.DB.Where("username = ?", cr.SenderName).Find(&senderUser).Count(&count)
	if count == 0 && senderUser.Phone != cr.SenderPhone {
		global.Log.Println("用户不存在直接建立")
		err := user_ser.UserService{}.CreateUser(cr.SenderName, "", ctype.PermissionUser, cr.SenderName, cr.SenderName)
		if err != nil {
			global.Log.Error(err)
			res.ResultFailWithMsg("用户创建失败!", c)
			return
		}
	}
	err := global.DB.Create(&item).Error
	if err != nil {
		res.ResultFailWithMsg("包裹创建失败!", c)
		return
	}
	res.ResultOkWithMsg("包裹创建成功!", c)
}
