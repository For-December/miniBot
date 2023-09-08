package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"testbot/apiUtils"
	"testbot/utils"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

var scheduleMap map[string]int

// Processor is a struct to process message
type Processor struct {
	Api openapi.OpenAPI
}

func init() {
	if scheduleMap == nil {
		scheduleMap = make(map[string]int) // 初始化 map
	}
}

// ProcessMessage is a function to process message
func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	toCreate := &dto.MessageToCreate{
		MsgID:   data.ID,
		Content: "默认回复" + message.Emoji(307),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}

	// 该艾特用户有待办任务要发送
	if userId := data.Author.ID; scheduleMap[userId] != 0 {
		// 设置日程
		task := message.ETLInput(input)
		utils.Info(fmt.Sprint(utils.IsTaskTxt(task)))
		if !utils.IsTaskTxt(task) {
			toCreate.Content = "您的待办事项格式有误，请修改后再次回复..." + message.Emoji(21) // 可爱
			p.sendReply(ctx, data.ChannelID, toCreate)
			return nil
		}

		toCreate.Content = "格式正确，开始设置待办事项，请稍后..." + message.Emoji(30) // 可爱
		p.sendReply(ctx, data.ChannelID, toCreate)
		delete(scheduleMap, data.Author.ID)
		params := utils.GetTaskParam(task)

		if !utils.IsRegistered(userId) {
			toCreate.Content = "用户未注册，为您自动注册中..."
			p.sendReply(ctx, data.ChannelID, toCreate)
			utils.CreateUsers(userId, data.Author.Username, "", "", "")
		}

		for i := range params["date"] {
			ok, info, task := utils.CreateTasks(userId,
				data.Author.Username, params["title"][i],
				params["context"][i],
				params["date"][i], "待办")
			toCreate.Content = info
			if !ok {
				if i > 0 {
					toCreate.Content += fmt.Sprintf("\n不过第 %v 个待办事项之前的所有任务都设置成功了...", i+1)
					toCreate.Content += message.Emoji(30) // 可爱
				} else {
					toCreate.Content += message.Emoji(38) // 敲打
				}
				p.sendReply(ctx, data.ChannelID, toCreate)
				return nil
			}

			utils.InfoF("设置用户 %v 的任务 %v : %v", data.Author.Username, i+1, params["title"][i])
			p.runTaskNoticeTimer(data.ChannelID,
				*task, true, "1921567337@qq.com")

		}

		toCreate.Content = "待办事项设置成功！" + message.Emoji(30)
		p.sendReply(ctx, data.ChannelID, toCreate)

		//for key, value := range params {
		//	utils.Info(key)
		//	utils.Info(value)
		//}

		return nil
	}

	// 进入到私信逻辑
	if cmd.Cmd == "dm" {
		p.dmHandler(data)
		return nil
	}

	guild, _ := p.Api.Guild(ctx, data.GuildID)
	channel, _ := p.Api.Channel(ctx, data.ChannelID)
	switch cmd.Cmd {
	case "打卡":
		ok := apiUtils.CreateForum(
			ChannelMap["话题区"],
			&apiUtils.ReqData{
				Title:   "实际测试",
				Content: "<html lang=\"en-US\"><body><a href=\"https://bot.q.qq.com/wiki\" title=\"QQ机器人文档Title\">QQ机器人文档</a>\n<ul><li>主动消息：发送消息时，未填msg_id字段的消息。</li><li>被动消息：发送消息时，填充了msg_id字段的消息。</li></ul></body></html>",
				Format:  2,
			})
		if ok {
			toCreate.Content = "打卡成功！"
			p.sendReply(ctx, data.ChannelID, toCreate)
		} else {
			toCreate.Content = "打卡失败！"
			p.sendReply(ctx, data.ChannelID, toCreate)
		}
	case "debug":
		toCreate.Content = ""
		for _, conversation := range WXChatData[data.Author.ID] {
			toCreate.Content += fmt.Sprint(conversation)
			toCreate.Content += "\r\n"
		}
		p.sendReply(ctx, data.ChannelID, toCreate)

	case "设置任务":
		if scheduleMap[data.Author.ID] == 0 {
			scheduleMap[data.Author.ID] += 1
			utils.Info("开始为用户" + data.Author.Username + "设置日程")
			toCreate.Content = `开始设置日程，可按格式设置多个日程。
请艾特我并按如下格式回复(可换行，email 字段可选)：`
			p.sendReply(ctx, data.ChannelID, toCreate)

			toCreate.Content = `
date: 2023/8/21-14:33:33
email: 123@xx
title: 标题
context: ` + "```内容```" + `
`
			p.sendReply(ctx, data.ChannelID, toCreate)

		}

	case "查询任务":
		tasksArray := utils.GetTasksById(data.Author.ID)
		toCreate.Content = "您所有的待办事项如下："
		p.sendReply(ctx, data.ChannelID, toCreate)
		for _, task := range tasksArray {
			toCreate.Content = fmt.Sprintf(`任务编号: %v
创建日期: %v
任务日期: %v
任务标题: %v
任务内容: %v`, task.TaskNum, task.CreatedDate, task.DueDate, task.Title, task.Description)
			p.sendReply(ctx, data.ChannelID, toCreate)

		}

	case "测试":
		toCreate.Content = guild.Name + "\n" + channel.Name
		//toCreate.Image = "https://qq-web.cdn-go.cn/im.qq.com_new/ca985481/img/product-tim.859a46a4.png"
		toCreate.Image = "https://test.fordece.cn/res/downloaded_image_1692458943.jpg"
		println(data.Author.ID)
		println(data.Author.Username)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "铯图":
		toCreate.Content = "图来啦~ " + message.Emoji(307)
		toCreate.Image = "https://test.fordece.cn/proxy"
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "翻译":
		switch channel.Name {
		case "霓虹":
			toCreate.Content = apiUtils.ToJp(cmd.Content)
			p.sendReply(ctx, data.ChannelID, toCreate)
		case "聊天室":
			toCreate.Content = apiUtils.ToZh(cmd.Content)
			p.sendReply(ctx, data.ChannelID, toCreate)
		case "阿妹你看":
			toCreate.Content = apiUtils.ToEn(cmd.Content)
			p.sendReply(ctx, data.ChannelID, toCreate)
		default:

		}

	case "翻译成中文":
		toCreate.Content = apiUtils.ToZh(cmd.Content)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "翻译成英文":
		toCreate.Content = apiUtils.ToEn(cmd.Content)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "翻译成日文":
		toCreate.Content = apiUtils.ToJp(cmd.Content)
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
		if isMatch, _ := regexp.MatchString("翻译\\S+?", cmd.Cmd); isMatch {
			toCreate.Content = `未找到命令，你可能想说：
> 翻译成日文
> 翻译成中文
> 翻译成英文
`
			p.sendReply(ctx, data.ChannelID, toCreate)
		} else {
			tempConversation := WXChatData[data.Author.ID]
			tempConversation = append(tempConversation,
				apiUtils.Conversation{Role: "user", Content: message.ETLInput(input)})

			reply := apiUtils.WXChat(tempConversation)
			WXChatData[data.Author.ID] = append(tempConversation, *reply)

			toCreate.Content = reply.Content
			p.sendReply(ctx, data.ChannelID, toCreate)
		}
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
	if err := p.Api.PutInteraction(context.Background(), interaction.ID, string(body)); err != nil {
		log.Println("apiUtils call putInteractionInlineSearch  error: ", err)
		return err
	}
	return nil
}

func (p Processor) ProcessDMMessage(data *dto.WSDirectMessageData) error {
	ctx := context.Background()

	_, err1 := p.Api.PostDirectMessage(ctx,
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
	dm, err := p.Api.CreateDirectMessage(
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
	_, err = p.Api.PostDirectMessage(
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
