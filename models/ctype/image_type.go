package ctype

import "github.com/goccy/go-json"

type ImageType int

const (
	Local ImageType = 1 //本地
)

func (r ImageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r ImageType) String() interface{} {
	var str string
	switch r {
	case Local:
		str = "本地"
	default:
		str = "未知"
	}
	return str
}
