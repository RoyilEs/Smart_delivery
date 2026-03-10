package config

type Admin struct {
	Username string `json:"username" yaml:"username"`
	Avatar   string `json:"avatar" yaml:"avatar"`
	Github   string `json:"github" yaml:"github"`
	BiliBili string `json:"bilibili" yaml:"bilibili"`
}
