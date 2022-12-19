package config

type redisKey struct {
	UserInfo                 string
	AccessToken              string
	YZM                      string
	UniDian                  string
	FriendsHelp              string
	SendPhoneByQnr           string
	Limit1S                  string
	InitTopicFlowLimit       string
	Lock                     string
	DuiBaShortUrl            string
	ActivityZeroIsNewUser    string
	OaAuth                   string
	BaiDu                    string
	BaiDuImageSearch         string
	BadgeImageCode           string
	CheckBusinessUser        string
	BusinessCarbonHotCity    string
	BlackList                string
	EventLimit               string
	UserCarbonRank           string
	UserCarbonClassify       string
	UserCarbonStep           string
	MessageLimitByTemplate   string
	MessageLimitByUser       string
	MessageSignUser          string
	MessageLimitTopicShow    string
	MessageLimitPlatformShow string
	MessageLimitCarbonPkShow string
	TopicRank                string
	YZM2B                    string
	GitlabHook               string
	CarbonPkRemindUser       string
}

var RedisKey = redisKey{
	UserInfo:                 "mp2c:userinfo:%d",                // 变量1:用户id 使用 fmt.Sprintf(RedisKey.UserInfo,"1")
	AccessToken:              "mp2c:access_token:%s:%s",         // 变量1:应用平台 变量2:应用id
	YZM:                      "yzm:",                            // 拼接用户id
	YZM2B:                    "yzm2b:",                          //验证码给b端，拼接用户id
	UniDian:                  "unidian:",                        // 拼接手机号
	FriendsHelp:              "friends_help:",                   // 拼接手机号
	SendPhoneByQnr:           "SendPhoneByQnr",                  //拼接手机号
	Limit1S:                  "Limit1S:",                        // 拼接行数名称
	InitTopicFlowLimit:       "mp2c:initTopicFlowlimit:%d",      //拼接用户id
	Lock:                     "mp2c:lock:",                      //redis分布式锁  拼接key
	DuiBaShortUrl:            "mp2c:duiba:shorturl:%s",          // 将对吧长链接存到redis中
	ActivityZeroIsNewUser:    "mp2c:activity:zero:isnewuser:%d", //0元拿活动记录是否新用户
	OaAuth:                   "mp2c:oaauth:%s",                  //微信网页授权
	BaiDu:                    "mp2c:baidu:",                     //用于百度接口
	BaiDuImageSearch:         "mp2c:baiduImageSearch",
	BadgeImageCode:           "mp2c:updateBadgeImage:",               //上传证书图片凭证
	CheckBusinessUser:        "mp2c:checkBusinessUser:%s",            //检测用户token是否有更新
	BusinessCarbonHotCity:    "business:carbon:hotcity",              //低碳场景中热门城市
	BlackList:                "mp2c:blacklist",                       //黑产白名单
	EventLimit:               "mp2c:eventlimit",                      //公益兑换限制
	UserCarbonRank:           "user_carbon_rank:%s",                  //排行榜
	UserCarbonClassify:       "user_carbon_classify",                 //碳分类
	UserCarbonStep:           "user_carbon_step:%s",                  //用户已经转化碳的步数
	MessageLimitByTemplate:   "mp2c:message_limit_by_template:%s_%s", //同一模板每人每天最多接收1条消息
	MessageLimitByUser:       "mp2c:message_limit_by_user:%s",        //每人每天最多收到2个不同类型模板消息
	MessageSignUser:          "mp2c:message_sign_user",               //小程序消息推送
	MessageLimitPlatformShow: "mp2c:message_limit_platform_show:%s",  //订阅消息每天弹出限制 平台普通
	MessageLimitTopicShow:    "mp2c:message_limit_topic_show:%s",     //订阅消息每天弹出限制 帖子
	MessageLimitCarbonPkShow: "mp2c:message_limit_carbon_pk_show:%s", //订阅消息每天弹出限制 打卡挑战

	TopicRank:          "mp2c:topic:rank",
	GitlabHook:         "mp2c:gitlab:hook:",
	CarbonPkRemindUser: "mp2c:carbon_pk_remind_user", //提醒用户池
}
