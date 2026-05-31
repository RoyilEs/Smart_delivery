package config

type FromMail struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	From     string `yaml:"from"`
	Password string `yaml:"password"`
}

type EmailConfig struct {
	FromMail FromMail `yaml:"from_mail"`
}
