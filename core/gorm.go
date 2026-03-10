package core

import (
	"Smart_delivery_locker/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitGorm() *gorm.DB {
	//未配置到mysql
	if global.Config.MySql.Host == "" {
		global.Log.Warning("未配置mysql, 请配置")
		return nil
	}
	dsn := global.Config.MySql.Dsn()

	var mySqlLogger logger.Interface
	if global.Config.System.Env == "debug" {
		//开发环境显示所有sql
		mySqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mySqlLogger = logger.Default.LogMode(logger.Error) //打印sql的错误
	}
	//gorm sql日志
	global.MySqlLog = logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mySqlLogger,
	})
	if err != nil {
		global.Log.Fatalf("[%s] mysql连接错误", dsn)
	}
	sqlDB, _ := db.DB()
	//TODO 设置mysql连接池 自定义配置 需要修改yaml与 conf—mysql
	sqlDB.SetMaxIdleConns(10)               //最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              //最多容量
	sqlDB.SetConnMaxLifetime(time.Hour * 4) //连接最大服用时间 不能超过mysql的wait_timeout
	return db
}
