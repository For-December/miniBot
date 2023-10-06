package main

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"strings"
	"testbot/controller"
	"testbot/logger"
)

// ReadyHandler 自定义 ReadyHandler 感知连接成功事件
func ReadyHandler() event.ReadyHandler {
	return func(event *dto.WSPayload, data *dto.WSReadyData) {
		//log.Println("ready event receive: ", data)
	}
}

func ErrorNotifyHandler() event.ErrorNotifyHandler {
	return func(err error) {
		//log.Println("error notify receive: ", err)
	}
}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content)) // 去掉@符号和首尾空格，同时全改小写
		return processor.ProcessMessage(input, data)             // 处理数据
	}
}

// GuildEventHandler 处理频道事件
func GuildEventHandler() event.GuildEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildData) error {
		//fmt.Println(data)
		return nil
	}
}

// ChannelEventHandler 处理子频道事件
func ChannelEventHandler() event.ChannelEventHandler {
	return func(event *dto.WSPayload, data *dto.WSChannelData) error {
		//fmt.Println(data)
		return nil
	}
}

// MemberEventHandler 处理成员变更事件
func MemberEventHandler() event.GuildMemberEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildMemberData) error {
		//fmt.Println(data)
		return nil
	}
}

// DirectMessageHandler 处理私信事件
func DirectMessageHandler() event.DirectMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSDirectMessageData) error {
		//fmt.Println(data.Content)

		return processor.ProcessDMMessage(data)
	}
}

// CreateMessageHandler 处理消息事件
func CreateMessageHandler() event.MessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageData) error {
		//fmt.Println(data)
		return nil
	}
}

// InteractionHandler 处理内联交互事件
func InteractionHandler() event.InteractionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		//fmt.Println(data)
		return processor.ProcessInlineSearch(data)
	}
}

// ThreadEventHandler 论坛主贴事件
func ThreadEventHandler() event.ThreadEventHandler {
	return func(event *dto.WSPayload, data *dto.WSThreadData) error {
		ctx := context.Background()

		logger.Debug(event.Type, data.ThreadInfo.DateTime)
		switch event.Type {
		case "FORUM_THREAD_UPDATE":
			logger.Warning("帖子被更新了！")
			logger.Info(string(event.RawMessage))
			processor.SendReply(ctx, controller.MainChannel.ID, &dto.MessageToCreate{Content: "<#" + data.GuildID + ">"})
			//apiUtils.ColorPicToChannel(controller.MainChannel.ID, ctx)
			//Embed: &dto.Embed{Title: "跳转",
			//	Description: "跳到测试的帖子",
			//	//Prompt:      "弹窗干扰",
			//	//Thumbnail:   dto.MessageEmbedThumbnail{URL: "https://mpqq.gtimg.cn/bot-wiki/online/images/introduce/Aspose.Words.a59f0707-65ac-4bec-8de6-d0d8efeb74d0.001.png"},
			//}})
			//}
		case "FORUM_THREAD_CREATE":
			logger.Warning("帖子被创建了！")
		default:
			logger.Warning("未知状况！")

		}

		return nil
	}
}
