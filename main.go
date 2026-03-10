package main

import (
	"Smart_delivery_locker/core"
	"Smart_delivery_locker/flag"
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/routers"
)

func main() {
	core.InitConf() // 配置文件读取
	global.Log = core.InitLogger()
	global.DB = core.InitGorm()

	//global.Redis = core.ConnectRedis() //TODO 这里win无法正常启用 在功能为使用前不开启

	//命令行参数绑定
	option := flag.Parse()
	if flag.IsWebStop(option) {
		flag.SwitchOption(option)
		return
	}

	router := routers.InitRouter()
	global.Log.Infof("启动成功:: %s", global.Config.System.Addr())
	err := router.Run(global.Config.System.Addr())
	if err != nil {
		global.Log.Error("启动失败")
	}

}
