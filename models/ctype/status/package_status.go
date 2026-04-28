package status

import "github.com/goccy/go-json"

type PackageStatus int

const (
	Created  PackageStatus = 0 // 已建立 待入柜
	Stored   PackageStatus = 1 // 已入柜
	PickedUp PackageStatus = 2 // 已取件
)

func (s PackageStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s PackageStatus) String() string {
	var str string
	switch s {
	case Created:
		str = "created"
	case Stored:
		str = "stored"
	case PickedUp:
		str = "picked_up"
	default:
		str = "nil"
	}
	return str
}
