package ctype

import "github.com/goccy/go-json"

type Size int

const (
	SizeLarge  Size = 3 // "大" 80 60 40
	SizeMedium Size = 2 // "中" 50 35 25
	SizeSmall  Size = 1 // "小" 30 20 15
	SizeXLarge Size = 4 // "超大" 120 80 60
)

func (s Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s Size) String() interface{} {
	var str string
	switch s {
	case SizeLarge:
		str = "large"
	case SizeMedium:
		str = "medium"
	case SizeSmall:
		str = "small"
	case SizeXLarge:
		str = "large"
	default:
		str = "未知"
	}
	return str
}
