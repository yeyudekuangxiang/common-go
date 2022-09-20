package config

//小程序订阅消息模版

type MessageTemplateId struct {
	ChangePoint          string
	OrderDeliver         string
	SignRemind           string
	TopicPass            string
	TopicCarefullyChosen string
}

var MessageTemplateIds = MessageTemplateId{
	ChangePoint:          "hQzUwkGMYqgNsKOad7RnIwwBpfkVfsuJvW6UqymwI8k", //积分到账提醒
	OrderDeliver:         "F8cnYtUUHc0blCkmvkXGIDtvr6v3tBqTE7LsTdvQJ4g", //订单发货提醒
	SignRemind:           "Oz2FJXoFAbjmE1s8s6GOqC2C5M9epKglDtz-3rjwJ6Q", //签到提醒
	TopicPass:            "0CNGli55ko4VBKdaUEiVbKWXEmu6tmnY",            //帖子审核通过提醒
	TopicCarefullyChosen: "G8_5XpSOsL0E0AmL0UfNZ3pwWRjU-PYdk40ptT_viFI", //帖子被精选通知
}
