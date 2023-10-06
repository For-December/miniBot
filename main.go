package main

import (
	"context"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"log"
	"reflect"
	"testbot/conf"
	"testbot/controller"
	"testbot/logger"
	"time"
)

// 消息处理器，持有 openapi 对象
var processor controller.Processor

func checkConf(value reflect.Value, lastFields ...string) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() == reflect.String {
			if field.IsZero() {
				resMsg := "未指定配置: "
				lastFields = append(lastFields, value.Type().Field(i).Name)
				for _, lastField := range lastFields {
					resMsg += lastField + "."
				}
				logger.Error(resMsg)
				return
			}
		} else {
			// 回溯！
			lastFields = append(lastFields, value.Type().Field(i).Name)
			checkConf(field, lastFields...)
			lastFields = lastFields[:len(lastFields)-1]

		}

	}
}
func init() {
	configValue := reflect.ValueOf(conf.Config)
	checkConf(configValue)

}

func main() {
	ctx := context.Background()
	botgo.SetLogger(logger.OverrideLogger{})

	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	if err := botToken.LoadFromConfig("config.yaml"); err != nil {
		log.Fatalln(err)
	}

	// 初始化 openapi，正式环境
	api := botgo.NewOpenAPI(botToken).WithTimeout(200 * time.Second)

	// 沙箱环境
	// apiUtils := botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second)

	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	processor = controller.Processor{Api: api}
	processor.InitAllTasks("机器人测试")
	processor.InitChanMap()

	// websocket.RegisterResumeSignal(syscall.SIGUSR1)
	// 根据不同的回调，生成 intents
	intent := websocket.RegisterHandlers(
		// at 机器人事件，目前是在这个事件处理中有逻辑，会回消息，其他的回调处理都只把数据打印出来，不做任何处理
		ATMessageEventHandler(),
		// 如果想要捕获到连接成功的事件，可以实现这个回调
		ReadyHandler(),
		// 连接关闭回调
		ErrorNotifyHandler(),
		// 频道事件
		GuildEventHandler(),
		// 成员事件
		MemberEventHandler(),
		// 子频道事件
		ChannelEventHandler(),
		// 私信，目前只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		DirectMessageHandler(),
		// 频道消息，只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		CreateMessageHandler(),
		// 互动事件
		InteractionHandler(),
		// 发帖事件
		ThreadEventHandler(),
	)
	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}
