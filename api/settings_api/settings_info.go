package settings_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models/res"
	CODE "Smart_delivery_locker/models/res/code"
	"github.com/gin-gonic/gin"
)

type SettingsUri struct {
	Name string `uri:"name"`
}

// SettingsInfoView 视图
func (SettingsApi) SettingsInfoView(c *gin.Context) {
	var cr SettingsUri
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}
	switch cr.Name {
	case "site":
		res.ResultOkWithData(global.Config.SiteInfo, c)
	case "jwt":
		res.ResultOkWithData(global.Config.Jwt, c)
	case "admin":
		res.ResultOkWithData(global.Config.Admin, c)
	default:
		res.ResultFailWithMsg("无对应配置", c)
	}

}
