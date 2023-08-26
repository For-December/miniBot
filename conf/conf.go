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

	// 日期格式
	DateLayout string `yaml:"dateLayout"`

	// 邮件通知
	Email struct {
		SenderEmail    string `yaml:"senderEmail"`
		SenderPassword string `yaml:"senderPassword"`
		SmtpServer     string `yaml:"smtpServer"`
	}
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
	Config.Appid = viper.GetString("appid")
	Config.Token = viper.GetString("token")

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
	Config.LogLevel = viper.GetString("logLevel")

	// 日期
	Config.DateLayout = viper.GetString("dateLayout")

	// 邮箱
	Config.Email.SmtpServer = viper.GetString("email.smtpServer")
	Config.Email.SenderEmail = viper.GetString("email.senderEmail")
	Config.Email.SenderPassword = viper.GetString("email.senderPassword")

}
