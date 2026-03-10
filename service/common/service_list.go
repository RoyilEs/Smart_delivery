package common

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"gorm.io/gorm"
)

type Option struct {
	models.PageInfo
	Debug bool
}

// ComList 分页查询
func ComList[T any](model T, option Option) (list []T, count int64, err error) {
	//引用sql日志
	DB := global.DB
	if option.Debug {
		DB = global.DB.Session(&gorm.Session{Logger: global.MySqlLog})
	}
	if option.Sort == "" {
		option.Sort = "created_at desc" //默认排序 按照时间往前排
	}

	if option.Limit == 0 {
		option.Limit = 10
	}

	//查询获得总数
	count = DB.Debug().Select("id").Find(&list).RowsAffected
	//分页 获得一共几页
	offset := (option.Page - 1) * option.Limit
	if offset < 0 {
		offset = 0
	}
	//分页后查询
	err = DB.Debug().Limit(option.Limit).Offset(offset).Order(option.Sort).Find(&list).Error
	return list, count, err
}
