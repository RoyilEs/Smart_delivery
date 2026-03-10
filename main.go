package main

import (
	"Smart_delivery_locker/core"
	"Smart_delivery_locker/global"
	"flag"
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
}
