package conf

import (
	"github.com/spf13/viper"
	"log"
	"reflect"
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
		SmtpServer     string `yaml:"smtpServer"`
		SmtpPort       string `yaml:"smtpPort"`
		SenderEmail    string `yaml:"senderEmail"`
		SenderPassword string `yaml:"senderPassword"`
	} `yaml:"email"`

	// AI
	AI struct {
		BaiduWX struct {
			ApiKey    string `yaml:"apiKey"`
			SecretKey string `yaml:"secretKey"`
		} `yaml:"baiduWX"`
	} `yaml:"ai"`

	Images struct {
		RandomApi string `yaml:"randomApi"`
	} `yaml:"images"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
		return // 自动退出
	}
	setConf(reflect.ValueOf(&Config))
}

// 利用反射，递归配置所有参数! recursion!!!
func setConf(value reflect.Value, lastFields ...string) {
	for i := 0; i < value.Elem().NumField(); i++ {
		field := value.Elem().Field(i)
		if field.Kind() == reflect.String {
			resKey := ""
			for _, lastField := range lastFields {
				resKey += lastField + "."
			}
			resKey += value.Type().Elem().Field(i).Name
			if tempParam := viper.GetString(resKey); tempParam != "" {
				field.Set(reflect.ValueOf(tempParam))
			}
		} else {
			// 回溯 (前进 => 处理 => 回退)
			lastFields = append(lastFields, value.Elem().Type().Field(i).Name)
			setConf(field.Addr(), lastFields...)
			lastFields = lastFields[:len(lastFields)-1]
		}
	}
}
