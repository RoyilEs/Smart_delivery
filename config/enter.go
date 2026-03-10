package config

// Config 实体化yaml信息
type Config struct {
	MySql    MySql    `yaml:"mysql"`
	Logger   Logger   `yaml:"logger"`
	System   System   `yaml:"system"`
	SiteInfo SiteInfo `yaml:"site_info"`
	Jwt      Jwt      `yaml:"jwt"`
	Upload   Upload   `yaml:"upload"`
	Redis    Redis    `yaml:"redis"`
	Admin    Admin    `yaml:"admin"`
}
