package controller

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"log"
	"testbot/conf"
	"testbot/dao"
	"testbot/logger"
	"testbot/utils"
	"time"
)

func (p Processor) setEmoji(ctx context.Context, channelID string, messageID string) {
	err := p.Api.CreateMessageReaction(
		ctx, channelID, messageID, dto.Emoji{
			ID:   "307",
			Type: 1,
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func (p Processor) setPins(ctx context.Context, channelID, msgID string) {
	_, err := p.Api.AddPins(ctx, channelID, msgID)
	if err != nil {
		log.Println(err)
	}
}

func (p Processor) setAnnounces(ctx context.Context, data *dto.WSATMessageData) {
	if _, err := p.Api.CreateChannelAnnounces(
		ctx, data.ChannelID,
		&dto.ChannelAnnouncesToCreate{MessageID: data.ID},
	); err != nil {
		log.Println(err)
	}
}

func (p Processor) SendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
	if _, err := p.Api.PostMessage(ctx, channelID, toCreate); err != nil {
		log.Println(err)
	}
}

func (p Processor) runTaskNoticeTimer(channelID string,
	tasks dao.Tasks,
	isEmail bool,
	email ...string) {

	// time.Parse 默认 UTC，这里指定本地地址（北京时间）
	dueDate, _ := time.ParseInLocation(conf.Config.DateLayout, tasks.DueDate, time.Local)
	// 未来的任务
	if dueDate.After(time.Now()) {
		logger.InfoF("一条未来的任务: {user: %v, num: %v}", tasks.Username, tasks.TaskNum)
		logger.Debug("dueDate: ", dueDate)
		logger.Debug("相差: ", dueDate.Sub(time.Now()))

		time.AfterFunc(dueDate.Sub(time.Now()), func() {
			// 定时艾特 + 邮箱提醒
			toCreate := &dto.MessageToCreate{
				Content: message.MentionUser(tasks.UserID) +
					"日程提醒：\r\n" +
					tasks.DueDate + "\r\n" +
					tasks.Title + "\r\n" +
					tasks.Description,
			}
			p.SendReply(context.Background(), channelID, toCreate)

			if isEmail {
				utils.SendEmail(email, "留意您的待办事项", toCreate.Content)
			}
		})
	}

}
