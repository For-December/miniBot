package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

// Processor is a struct to process message
type Processor struct {
	api openapi.OpenAPI
}

// ProcessMessage is a function to process message
func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	toCreate := &dto.MessageToCreate{
		Content: "默认回复" + message.Emoji(307),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}

	// 进入到私信逻辑
	if cmd.Cmd == "dm" {
		p.dmHandler(data)
		return nil
	}
	guild, _ := p.api.Guild(ctx, data.GuildID)
	channel, _ := p.api.Channel(ctx, data.ChannelID)
	switch cmd.Cmd {
	case "测试":
		toCreate.Content = guild.Name + "\n" + channel.Name
		//toCreate.Image = "https://qq-web.cdn-go.cn/im.qq.com_new/ca985481/img/product-tim.859a46a4.png"
		toCreate.Image = "https://test.fordece.cn/res/downloaded_image_1692458943.jpg"
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "铯图":
		toCreate.Content = "图来啦~ " + message.Emoji(307)
		log.Println(colorYellow + "***************************" + colorReset)
		toCreate.Image = "https://test.fordece.cn/proxy"
		//toCreate.Image = "https://test.fordece.cn/res/downloaded_image_1692458943.jpg"
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "翻译":
		switch channel.Name {
		case "霓虹":
			toCreate.Content = toJp(cmd.Content)
			p.sendReply(ctx, data.ChannelID, toCreate)
		case "聊天室":
			toCreate.Content = toZh(cmd.Content)
			p.sendReply(ctx, data.ChannelID, toCreate)
		case "阿妹你看":
			toCreate.Content = toEn(cmd.Content)
			p.sendReply(ctx, data.ChannelID, toCreate)
		default:

		}

	case "翻译成中文":
		toCreate.Content = toZh(cmd.Content)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "翻译成英文":
		toCreate.Content = toEn(cmd.Content)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "翻译成日文":
		toCreate.Content = toJp(cmd.Content)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "hi":
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "time":
		toCreate.Content = genReplyContent(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "ark":
		toCreate.Ark = genReplyArk(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "公告":
		p.setAnnounces(ctx, data)
	case "pin":
		if data.MessageReference != nil {
			p.setPins(ctx, data.ChannelID, data.MessageReference.MessageID)
		}
	case "emoji":
		if data.MessageReference != nil {
			p.setEmoji(ctx, data.ChannelID, data.MessageReference.MessageID)
		}
	default:
	}

	return nil
}

// ProcessInlineSearch is a function to process inline search
func (p Processor) ProcessInlineSearch(interaction *dto.WSInteractionData) error {
	if interaction.Data.Type != dto.InteractionDataTypeChatSearch {
		return fmt.Errorf("interaction data type not chat search")
	}
	search := &dto.SearchInputResolved{}
	if err := json.Unmarshal(interaction.Data.Resolved, search); err != nil {
		log.Println(err)
		return err
	}
	if search.Keyword != "test" {
		return fmt.Errorf("resolved search key not allowed")
	}
	searchRsp := &dto.SearchRsp{
		Layouts: []dto.SearchLayout{
			{
				LayoutType: 0,
				ActionType: 0,
				Title:      "内联搜索",
				Records: []dto.SearchRecord{
					{
						Cover: "https://pub.idqqimg.com/pc/misc/files/20211208/311cfc87ce394c62b7c9f0508658cf25.png",
						Title: "内联搜索标题",
						Tips:  "内联搜索 tips",
						URL:   "https://www.qq.com",
					},
				},
			},
		},
	}
	body, _ := json.Marshal(searchRsp)
	if err := p.api.PutInteraction(context.Background(), interaction.ID, string(body)); err != nil {
		log.Println("api call putInteractionInlineSearch  error: ", err)
		return err
	}
	return nil
}

func (p Processor) ProcessDMMessage(data *dto.WSDirectMessageData) error {
	ctx := context.Background()

	_, err1 := p.api.PostDirectMessage(ctx,
		&dto.DirectMessage{GuildID: data.GuildID},
		&dto.MessageToCreate{
			Content: "私信消息,不知道该怎么回复你",
			MsgID:   data.ID,
		})
	if err1 != nil {
		log.Fatalln("调用 PostDirectMessage 接口失败, err = ", err1)
	}
	return nil
}

func (p Processor) dmHandler(data *dto.WSATMessageData) {
	dm, err := p.api.CreateDirectMessage(
		context.Background(), &dto.DirectMessageToCreate{
			SourceGuildID: data.GuildID,
			RecipientID:   data.Author.ID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	toCreate := &dto.MessageToCreate{
		Content: "默认私信回复",
	}
	_, err = p.api.PostDirectMessage(
		context.Background(), dm, toCreate,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func genReplyContent(data *dto.WSATMessageData) string {
	var tpl = `你好：%s
在子频道 %s 收到消息。
收到的消息发送时时间为：%s
当前本地时间为：%s

消息来自：%s
`

	msgTime, _ := data.Timestamp.Time()
	return fmt.Sprintf(
		tpl,
		message.MentionUser(data.Author.ID),
		message.MentionChannel(data.ChannelID),
		msgTime, time.Now().Format(time.RFC3339),
		getIP(),
	)
}

func genReplyArk(data *dto.WSATMessageData) *dto.Ark {
	return &dto.Ark{
		TemplateID: 23,
		KV: []*dto.ArkKV{
			{
				Key:   "#DESC#",
				Value: "这是 ark 的描述信息",
			},
			{
				Key:   "#PROMPT#",
				Value: "这是 ark 的摘要信息",
			},
			{
				Key: "#LIST#",
				Obj: []*dto.ArkObj{
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这里展示的是 23 号模板",
							},
						},
					},
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这是 ark 的列表项名称",
							},
							{
								Key:   "link",
								Value: "https://www.qq.com",
							},
						},
					},
				},
			},
		},
	}
}
