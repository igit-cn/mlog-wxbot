package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Conf *Config

type Config struct {
	MySqlUrl string `yaml:"MySqlUrl"` // 数据库连接地址
	ShowSql  bool   `yaml:"ShowSql"`  // 是否显示日志

	IgnoreGzhConfig string `yaml:"IgnoreGzhConfig"` // 需要过滤的公众号配置文件，文件中一行一个公众号名称

	// 百度ai
	BaiduAi struct {
		ApiKey    string `yaml:"ApiKey"`
		SecretKey string `yaml:"SecretKey"`
	} `yaml:"BaiduAi"`

	PublishToken string `yaml:"PublishToken"` // 发表文章所需要的token
	PublishApi   string `yaml:"PublishApi"`   // 文章发布api

	// 阿里云oss配置
	AliyunOss struct {
		Host         string `yaml:"Host"`
		Bucket       string `yaml:"Bucket"`
		Endpoint     string `yaml:"Endpoint"`
		AccessId     string `yaml:"AccessId"`
		AccessSecret string `yaml:"AccessSecret"`
	} `yaml:"AliyunOss"`
}

func InitConfig(filename string) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Error(err)
		return
	}

	Conf = &Config{}
	err = yaml.Unmarshal(yamlFile, Conf)
	if err != nil {
		logrus.Error(err)
	}
}
