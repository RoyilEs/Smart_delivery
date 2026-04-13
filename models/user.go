package models

import (
	"Smart_delivery_locker/models/ctype"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	Username   string     `json:"username" gorm:"type:varchar(20);not null;unique"` //用户名
	Nickname   string     `json:"nickname" gorm:"type:varchar(20);"`                //昵称
	Email      string     `json:"email" gorm:"type:varchar(255);"`                  //邮箱
	Phone      string     `json:"phone" gorm:"type:varchar(20);"`                   //手机号
	Password   string     `json:"password" gorm:"type:varchar(255);not null"`       //密码
	Permission ctype.Role `json:"permission" gorm:"type:int(10);not null"`          //权限
	Avatar     string     `json:"avatar" gorm:"type:varchar(255);"`                 //头像
	Token      string     `json:"token" gorm:"size:64"`                             //其他平台的唯一地址
	Status     string     `json:"status" gorm:"type:varchar(20);"`                  //状态
}
