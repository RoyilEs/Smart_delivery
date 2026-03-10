package user_ser

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype"
	"Smart_delivery_locker/service/redis_ser"
	"Smart_delivery_locker/utils/jwts"
	"Smart_delivery_locker/utils/pwd"
	"errors"
	"time"
)

type UserService struct {
}

// Logout 是否登出
func (UserService) Logout(claims *jwts.CustomClaims, token string) error {
	//过期时间 需要计算过期时间 距离现在的过期时间
	exp := claims.ExpiresAt
	now := time.Now()
	// 截至时间
	diff := exp.Time.Sub(now)
	return redis_ser.Logout(token, diff)
}

const AVATAR = "/uploads/avatar/头像.png"

func (UserService) CreateUser(userName, password string, permission ctype.Role, email string, phone string) error {
	//判断逻辑结构
	//判断用户名是否存在
	var userModel models.User
	err := global.DB.Take(&userModel, "username = ?", userName).Error
	if err == nil {
		//存在
		return errors.New("用户名已存在,请重新输入")
	}
	//TODO 正则密码强度
	//加密密码 hash
	hashPassword := pwd.HashPassword(password)

	//头像问题 1.默认 2.随机选择

	//入库
	err = global.DB.Create(&models.User{
		Username:   userName,
		Email:      email,
		Phone:      phone,
		Password:   hashPassword,
		Permission: permission,
		Avatar:     AVATAR,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
