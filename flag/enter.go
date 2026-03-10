package flag

import (
	FLAG "flag"
	"github.com/fatih/structs"
)

type Option struct {
	DB   bool
	User string // -u admin -u user
}

// Parse 解析命令行参数
func Parse() Option {
	db := FLAG.Bool("db", false, "初始化数据库")
	user := FLAG.String("u", "", "创建用户")
	//解析命令写入注册的flag中
	FLAG.Parse()
	return Option{
		DB:   *db,
		User: *user,
	}
}

// IsWebStop 是否停止web项目
func IsWebStop(option Option) (f bool) {
	maps := structs.Map(&option)
	for _, val := range maps {
		switch v := val.(type) {
		case string:
			if v != "" {
				f = true
			}
		case bool:
			if v == true {
				f = true
			}
		}
	}
	return
}

// SwitchOption 根据命令执行不同的函数
func SwitchOption(option Option) {
	if option.DB {
		Makemigrations()
		return
	}
	if option.User == "admin" || option.User == "user" {
		CreateUser(option.User)
		return
	}
	//FLAG.Usage()
}
