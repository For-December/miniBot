package controller

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
	"testbot/apiUtils"
	"testbot/utils"
)

var WXChatData = make(map[string][]apiUtils.Conversation)

var ChannelMap = make(map[string]string)

func (p Processor) InitChanMap() {
	ctx := context.Background()
	guilds, err := p.Api.MeGuilds(ctx, &dto.GuildPager{Limit: "1"})
	if err != nil {
		return
	}

	channels, err := p.Api.Channels(ctx, guilds[0].ID)
	if err != nil {
		utils.Error("出错了: ", err)
	}

	for _, channel := range channels {
		ChannelMap[channel.Name] = channel.ID
	}
}
func (p Processor) InitAllTasks(channelName string) {
	utils.InfoF("为子频道 %v 配置通知...", channelName)
	ctx := context.Background()
	guilds, err := p.Api.MeGuilds(ctx, &dto.GuildPager{Limit: "1"})
	if err != nil {
		return
	}

	channel := func(channelName string) *dto.Channel {
		channels, err := p.Api.Channels(ctx, guilds[0].ID)
		if err != nil {
			utils.Error("出错了: ", err)
		}

		for _, channel := range channels {
			if channel.Name == channelName {
				return channel
			}
		}
		utils.ErrorF("未找到子频道 %v！", channelName)
		return nil
	}(channelName)

	tasksArray := utils.GetTasks()
	utils.InfoF("为子频道 %v 初始化任务中...", channelName)

	for _, tasks := range tasksArray {
		p.runTaskNoticeTimer(channel.ID, tasks, true,
			"1921567337@qq.com")
	}
}
