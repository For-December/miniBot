package apiUtils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"io"
	"net/http"
	"testbot/conf"
	"testbot/logger"
	"time"
)

type ReqData struct {
	MsgID   string `json:"msg_id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Format  int    `json:"format,omitempty"`
}

type respData struct {
	TaskId     string `json:"task_id,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type postsReply struct {
	Threads []struct {
		GuildId    string `json:"guild_id"`
		ChannelId  string `json:"channel_id"`
		AuthorId   string `json:"author_id"`
		ThreadInfo struct {
			ThreadId string    `json:"thread_id"`
			Title    string    `json:"title"`
			Content  string    `json:"content"`
			DateTime time.Time `json:"date_time"`
		} `json:"thread_info"`
	} `json:"threads"`
	IsFinish int `json:"is_finish"`
}

func CreateForum(channelID string, msg *ReqData) bool {
	payload, _ := json.Marshal(msg)
	logger.Debug(msg.Content)
	url := fmt.Sprintf(
		"https://api.sgroup.qq.com/channels/%s/threads", channelID)

	logger.Debug(string(payload))
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(req.Body)
	req.Header.Set("authorization",
		fmt.Sprintf("Bot %s.%s",
			conf.Config.Appid, conf.Config.Token))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false

}

func SendPicToChannelMsg(
	channelID string,
	qrContent []byte,
	data map[string]string,
	ctx context.Context) ([]byte, error) {
	resp, err := resty.New().R().SetContext(ctx).SetAuthScheme("Bot").
		SetAuthToken(conf.Config.Appid+"."+conf.Config.Token).
		SetFormData(data).
		SetFileReader("file_image", "qrcode.png", bytes.NewReader(qrContent)).
		SetContentLength(true).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		Post(fmt.Sprintf("%s://%s%s", "https", "api.sgroup.qq.com", "/channels/{channel_id}/messages"))
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func ColorPicToChannel(channelId string, ctx context.Context) {
	// 随机api
	resp, err := http.Get(conf.Config.Images.RandomApi) // https://test.fordece.cn/proxy
	if err != nil {
		logger.Error(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(resp.Body)

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}

	msg, err := SendPicToChannelMsg(channelId, imgBytes, map[string]string{
		"content": "图来啦" + message.Emoji(307),
	}, ctx)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(string(msg))
}
