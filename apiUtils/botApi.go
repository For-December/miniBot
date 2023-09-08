package apiUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testbot/conf"
	"testbot/utils"
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
	utils.Debug(msg.Content)
	url := fmt.Sprintf(
		"https://api.sgroup.qq.com/channels/%s/threads", channelID)

	utils.Debug(string(payload))
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Error(err)
		}
	}(req.Body)
	req.Header.Set("authorization",
		fmt.Sprintf("Bot %s.%s",
			conf.Config.Appid, conf.Config.Token))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.Error(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Error(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false

}
