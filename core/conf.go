package core

import (
	"Smart_delivery_locker/config"
	"Smart_delivery_locker/global"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/fs"
	"io/ioutil"
	"log"
)

const ConfigFile = "application.yaml"

// InitConf conf初始化读取配置文件
func InitConf() {

	c := &config.Config{}
	yamlConf, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get yamlConf error: %v", err))
	}
	//读取配置文件
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("yaml.Unmarshal error:%v", err)
	}
	log.Println("config yamlFile InitConf success")
	//设置到global-config
	global.Config = c
}

func SetYaml() (err error) {
	marshal, err := yaml.Marshal(global.Config)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(ConfigFile, marshal, fs.ModePerm)
	if err != nil {
		return
	}
	global.Log.Info("config yamlFile SetYaml success d=====(￣▽￣*)b")
	return nil
}
