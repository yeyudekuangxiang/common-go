package message

type MessageTemplateId struct {
	ChangePoint          string
	OrderDeliver         string
	SignRemind           string
	TopicPass            string
	TopicCarefullyChosen string
	PunchClockRemind     string
	QuizRemind           string
	ChargeOrder          string
}

type MessageSendMixCount struct {
	ChangePoint          float64
	OrderDeliver         float64
	SignRemind           float64
	TopicPass            float64
	TopicCarefullyChosen float64
	PunchClockRemind     float64
	QuizRemind           float64
	ChargeOrder          float64
}

//模版id

var MessageTemplateIds = MessageTemplateId{
	ChangePoint:          "hQzUwkGMYqgNsKOad7RnIwwBpfkVfsuJvW6UqymwI8k", //积分到账提醒
	OrderDeliver:         "F8cnYtUUHc0blCkmvkXGIDtvr6v3tBqTE7LsTdvQJ4g", //订单发货提醒
	SignRemind:           "Oz2FJXoFAbjmE1s8s6GOqC2C5M9epKglDtz-3rjwJ6Q", //签到提醒
	TopicPass:            "0CNGli55ko4VBKdaUEiVbKWXEmu6tmnY",            //帖子审核通过提醒
	TopicCarefullyChosen: "G8_5XpSOsL0E0AmL0UfNZ3pwWRjU-PYdk40ptT_viFI", //帖子被精选通知
	PunchClockRemind:     "lxV--pg9udJn_iJL0txG00rRqZiwJYLGOPm9g2dEo4Q", //打开提醒
	QuizRemind:           "8AHI8Iqd-HzSyYQIeNpZ-L2tKeJm8IkEneoBcFL0OkA", //答题闯关提醒
	ChargeOrder:          "atTU5emLvZAF88DFQCF4bT5F3ihH-yS0IcdhWF9pAc4", //绿喵充电
}

//模版每天最多发送条数

var MessageSendMixCounts = MessageSendMixCount{
	ChangePoint:          1, //积分到账提醒
	OrderDeliver:         1, //订单发货提醒
	SignRemind:           1, //签到提醒
	TopicPass:            1, //帖子审核通过提醒
	TopicCarefullyChosen: 1, //帖子被精选通知
	PunchClockRemind:     1, //打卡提醒
	QuizRemind:           1, //答题闯关提醒
	ChargeOrder:          20,
}

//模版跳转名称

var MessageTemplateName = MessageTemplateId{
	ChangePoint:          "积分变更提醒",
	OrderDeliver:         "订单发货提醒",
	SignRemind:           "签到提醒",
	TopicPass:            "帖子审核通过提醒", //
	TopicCarefullyChosen: "帖子被精选通知",  //
	PunchClockRemind:     "打卡提醒",
	QuizRemind:           "答题挑战提醒",
	ChargeOrder:          "充电结束订单提醒",
}

//模版跳转路径

var MessageJumpUrls = MessageTemplateId{
	ChangePoint:          "/pages/my_info/integral/index",                                                                //积分到账提醒
	OrderDeliver:         "/pages/duiba_v2/share/index?activityId=duiba_order",                                           //订单发货提醒
	SignRemind:           "/pages/duiba_v2/share/index?activityId=duiba_sign_Subnews&cid=1055&bind=true&hideShare=false", //签到提醒
	TopicPass:            "index",                                                                                        //帖子审核通过提醒
	TopicCarefullyChosen: "index",                                                                                        //帖子被精选通知
	PunchClockRemind:     "/pages/activity/punch/start/index",
	QuizRemind:           "/pages/answer_game/index",
	ChargeOrder:          "/pages/scene/charge/order/index?id=%s",
}
