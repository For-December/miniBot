package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Reply struct {
	From        string              `json:"from"`
	To          string              `json:"to"`
	TransResult []map[string]string `json:"trans_result"`
	ErrorCode   string              `json:"error_code"`
	ErrorMsg    string              `json:"error_msg"`
}

var appid string
var key string
var salt string
var apiUrl string

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
		return // 自动退出
	}

	appid = viper.GetString("baiduTrans.appid")
	key = viper.GetString("baiduTrans.key")
	salt = viper.GetString("baiduTrans.salt")
	apiUrl = viper.GetString("baiduTrans.apiUrl")
}

func toEn(q string) string {
	return call(q, "en")
}
func toJp(q string) string {
	return call(q, "jp")
}
func toZh(q string) string {
	return call(q, "zh")
}

func call(q string, dstLanguage string) (res string) {

	// 创建 MD5 哈希对象
	hashes := md5.New()

	// 将字符串写入哈希对象
	hashes.Write([]byte(appid + q + salt + key))

	// 计算 MD5 值
	sign := hex.EncodeToString(hashes.Sum(nil))

	//fmt.Printf("md5String: %v\n", sign)

	// 配置 post 参数
	data := url.Values{}
	data.Set("q", q)
	data.Set("from", "auto")
	data.Set("to", dstLanguage)
	data.Set("appid", appid)
	data.Set("salt", salt)
	data.Set("sign", sign)

	resp, _ := http.Post(apiUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()))

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("responseBody: %v\n", string(responseBody))

	var reply Reply
	err := json.Unmarshal(responseBody, &reply)
	if err != nil {
		log.Println(err)
		return
	}
	if reply.ErrorCode != "" {
		log.Println("请求api出错")
		res = "出错了，错误信息：" + reply.ErrorMsg
		return
	}

	res = reply.TransResult[0]["dst"]
	return

}
