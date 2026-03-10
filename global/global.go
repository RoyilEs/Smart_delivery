package global

import (
	"Smart_delivery_locker/config"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Config         *config.Config
	DB             *gorm.DB
	Log            *logrus.Logger
	MySqlLog       logger.Interface
	WhiteImageList = []string{
		"jpg",
		"png",
		"jpeg",
		"ico",
		"tiff",
		"gif",
		"svg",
		"webp",
		"bmp",
	} //图片白名单
	Redis *redis.Client
)
