package action

import (
	"github.com/goccy/go-json"
)

type PackageAction int

const (
	Created  PackageAction = 1 // 建立包裹
	Stored   PackageAction = 2 // 入柜成功
	PickedUp PackageAction = 3 // 已取件
	Updated  PackageAction = 4 // 后台编辑包裹
)

func (s PackageAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s PackageAction) String() string {
	var str string
	switch s {
	case Created:
		str = "created"
	case Stored:
		str = "stored"
	case PickedUp:
		str = "picked_up"
	case Updated:
		str = "updated"
	default:
		str = "nil"
	}
	return str
}
