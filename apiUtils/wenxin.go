package apiUtils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testbot/conf"
	"testbot/logger"
)

var apiKey string
var secretKey string

func init() {
	apiKey = conf.Config.AI.BaiduWX.ApiKey
	secretKey = conf.Config.AI.BaiduWX.SecretKey
}

type Conversation struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type WXReply struct {
	Id               string `json:"id"`
	Object           string `json:"object"`
	Created          int    `json:"created"`
	Result           string `json:"result"`
	IsTruncated      bool   `json:"is_truncated"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func WXChat(conversation []Conversation) *Conversation {
	conversationBytes, _ := json.Marshal(conversation)
	logger.Debug(string(conversationBytes))
	url := "https://aip.baidubce.com/" +
		"rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant?" +
		"access_token=" + getAccessToken()
	payload := strings.NewReader(fmt.Sprintf(`{"messages":%s}`, string(conversationBytes)))
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logger.Error(err)
		return nil
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}

	var reply WXReply
	err = json.Unmarshal(body, &reply)
	if err != nil {
		logger.Error(err)
		return nil
	}
	fmt.Println(reply.Result)

	if reply.ErrorCode != "" {
		return &Conversation{Role: "assistant", Content: "出错了: " + reply.ErrorMsg}
	}
	return &Conversation{Role: "assistant", Content: reply.Result}
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func getAccessToken() string {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials"+
		"&client_id=%s"+
		"&client_secret=%s", apiKey, secretKey)
	resp, err := http.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader(postData))
	if err != nil {
		logger.Error(err)
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return ""
	}
	accessTokenObj := map[string]any{}
	err = json.Unmarshal(body, &accessTokenObj)
	//utils.Debug(accessTokenObj)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return accessTokenObj["access_token"].(string)
}
