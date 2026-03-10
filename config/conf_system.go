package config

import "fmt"

type System struct {
	Host string `yaml:"host"`
	Post int    `yaml:"post"`
	Env  string `yaml:"env"`
}

// Addr 拼接路由
func (s *System) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Post)
}
