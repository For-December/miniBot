package controller

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"testbot/conf"
	"testbot/utils"
	"time"
)

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
		// time.Parse 默认 UTC
		dueDate, _ := time.ParseInLocation(conf.Config.DateLayout, tasks.DueDate, time.Local)
		// 未来的任务
		if dueDate.After(time.Now()) {
			utils.InfoF("一条未来的任务: {user: %v, num: %v}", tasks.Username, tasks.TaskNum)
			utils.Info("dueDate: ", dueDate)
			utils.Info("相差: ", dueDate.Sub(time.Now()))
			time.AfterFunc(dueDate.Sub(time.Now()), func() {
				// 定时艾特
				toCreate := &dto.MessageToCreate{
					Content: message.MentionUser(tasks.UserID) +
						"日程提醒：\n" +
						tasks.DueDate + "\n" +
						tasks.Title + "\n" +
						tasks.Description,
				}
				p.sendReply(ctx, channel.ID, toCreate)
			})
		}

	}
}
