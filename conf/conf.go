package conf

import (
	"github.com/spf13/viper"
	"log"
)

var Config struct {
	// botToken
	Appid string `yaml:"appid"`
	Token string `yaml:"token"`

	// 百度翻译api
	BaiduTrans struct {
		Appid  string `yaml:"appid"`
		Key    string `yaml:"key"`
		Salt   string `yaml:"salt"`
		ApiUrl string `yaml:"apiUrl"`
	} `yaml:"baiduTrans"`

	// 数据库配置，mysqlConf
	Mysql struct {
		DatabaseName string `yaml:"databaseName"`
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
	} `yaml:"mysql"`

	// 日志等级
	LogLevel string `yaml:"logLevel"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
		return // 自动退出
	}
	// bot
	// 已经在main函数中获取完成

	// 百度
	Config.BaiduTrans.Appid = viper.GetString("baiduTrans.appid")
	Config.BaiduTrans.Key = viper.GetString("baiduTrans.key")
	Config.BaiduTrans.Salt = viper.GetString("baiduTrans.salt")
	Config.BaiduTrans.ApiUrl = viper.GetString("baiduTrans.apiUrl")

	// mysql
	Config.Mysql.DatabaseName = viper.GetString("mysql.databaseName")
	Config.Mysql.Host = viper.GetString("mysql.host")
	Config.Mysql.Port = viper.GetString("mysql.port")
	Config.Mysql.Username = viper.GetString("mysql.username")
	Config.Mysql.Password = viper.GetString("mysql.password")

	// 日志
	Config.LogLevel = "debug"

}
