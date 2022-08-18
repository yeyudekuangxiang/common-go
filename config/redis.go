package config

type redisKey struct {
	UserInfo              string
	AccessToken           string
	YZM                   string
	UniDian               string
	FriendsHelp           string
	Limit1S               string
	InitTopicFlowLimit    string
	Lock                  string
	DuiBaShortUrl         string
	ActivityZeroIsNewUser string
	OaAuth                string
	BaiDu                 string
	BadgeImageCode        string
	CheckBusinessUser     string
	BlackList             string
}

var RedisKey = redisKey{
	UserInfo:              "mp2c:userinfo:%d",                // 变量1:用户id 使用 fmt.Sprintf(RedisKey.UserInfo,"1")
	AccessToken:           "mp2c:access_token:%s:%s",         // 变量1:应用平台 变量2:应用id
	YZM:                   "yzm:",                            // 拼接用户id
	UniDian:               "unidian:",                        // 拼接手机号
	FriendsHelp:           "friends_help:",                   // 拼接手机号
	Limit1S:               "Limit1S:",                        // 拼接行数名称
	InitTopicFlowLimit:    "mp2c:initTopicFlowlimit:%d",      //拼接用户id
	Lock:                  "mp2c:lock:",                      //redis分布式锁  拼接key
	DuiBaShortUrl:         "mp2c:duiba:shorturl:%s",          // 将对吧长链接存到redis中
	ActivityZeroIsNewUser: "mp2c:activity:zero:isnewuser:%d", //0元拿活动记录是否新用户
	OaAuth:                "mp2c:oaauth:%s",                  //微信网页授权
	BaiDu:                 "mp2c:baidu:",                     //用于百度接口
	BadgeImageCode:        "mp2c:updateBadgeImage:",          //上传证书图片凭证
	CheckBusinessUser:     "mp2c:checkBusinessUser:%s",       //检测用户token是否有更新
	BlackList:             "mp2c:blacklist",                  //黑产白名单
}
