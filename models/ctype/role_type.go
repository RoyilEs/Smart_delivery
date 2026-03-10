package ctype

import "github.com/goccy/go-json"

type Role int

const (
	PermissionAdmin       Role = 1 //管理员
	PermissionUser        Role = 2 //用户
	PermissionCourier     Role = 3 //快递员
	PermissionDisableUser Role = 4 //被禁用
)

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r Role) String() interface{} {
	var str string
	switch r {
	case PermissionAdmin:
		str = "管理员"
	case PermissionUser:
		str = "用户"
	case PermissionCourier:
		str = "快递员"
	case PermissionDisableUser:
		str = "被禁用"
	default:
		str = "未知"
	}
	return str
}
