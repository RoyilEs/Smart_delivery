package flag

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models/ctype"
	"fmt"
)

var (
	userName   string
	password   string
	rePassword string
	email      string
)

func CreateUser(permissions string) {

	fmt.Printf("请输入用户名:")
	fmt.Scan(&userName)
	fmt.Printf("请输入邮箱:")
	fmt.Scan(&email)
	fmt.Printf("请输入密码:")
	fmt.Scan(&password)
	fmt.Printf("请确定密码:")
	fmt.Scan(&rePassword)

	fmt.Println(toString())

	// 角色判断
	role := ctype.PermissionCourier
	if permissions == "admin" {
		role = ctype.PermissionAdmin
	}

	if password != rePassword {
		global.Log.Error("两次密码不一致,请重新输入")
		return
	}
}

func toString() string {
	return fmt.Sprintf("用户名:%s \\密码:%s \\确定密码:%s \\邮箱:%s", userName, password, rePassword, email)
}
