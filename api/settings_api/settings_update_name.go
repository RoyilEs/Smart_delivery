package settings_api

import (
	"Smart_delivery_locker/core"
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models/res"
	CODE "Smart_delivery_locker/models/res/code"
	"github.com/gin-gonic/gin"
)

func (SettingsApi) SettingsUpdateNameView(c *gin.Context) {
	name := c.Param("name")
	switch name {
	case "basic":
		updateConfig(c, &global.Config.Basic, "basic")
	case "pickup":
		updateConfig(c, &global.Config.Pickup, "pickup")
	}
}

func updateConfig[T any](c *gin.Context, configPtr *T, configName string) {
	var cr T
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}
	*configPtr = cr
	err = core.SetYaml()
	if err != nil {
		global.Log.Error(err)
		res.ResultFailWithMsg(err.Error(), c)
		return
	}
	res.ResultOkWithData(*configPtr, c)
}
