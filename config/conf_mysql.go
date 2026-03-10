package config

import "strconv"

type MySql struct {
	DataBase   string `yaml:"dataBase"`
	UserName   string `yaml:"userName"`
	PassWord   string `yaml:"passWord"`
	Port       int    `yaml:"port"`
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	LogLevel   string `yaml:"log_Level"` //日志等级 debug输出全部sql dev, release
}

func (m *MySql) Dsn() string {
	return m.UserName + ":" + m.PassWord + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.DataBase + "?charset=utf8&parseTime=True"
}
