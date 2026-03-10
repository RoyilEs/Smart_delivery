package config

type Upload struct {
	Size       int    `json:"size" yaml:"size"`               //图片上传的大小限制
	SizeAvatar int    `json:"size_avatar" yaml:"size_avatar"` //头像图片上传的大小限制
	Path       string `json:"path" yaml:"path"`               //图片上传的路径
	PathAvatar string `json:"path_avatar" yaml:"path_avatar"` //头像上传的路径
}
