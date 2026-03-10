package settings_api

import (
	"Smart_delivery_locker/config"
	"Smart_delivery_locker/core"
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models/res"
	CODE "Smart_delivery_locker/models/res/code"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (SettingsApi) SettingsJwtUpdateView(c *gin.Context) {
	var cr config.Jwt
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}
	fmt.Println("1====:", global.Config.Jwt)
	//修改
	global.Config.Jwt = cr
	err = core.SetYaml()
	if err != nil {
		global.Log.Error(err)
		res.ResultFailWithMsg(err.Error(), c)
		return
	}
	res.ResultOkWithMsg("修改成功d=====(￣▽￣*)b", c)
	fmt.Println("2====:", global.Config.Jwt)
}
