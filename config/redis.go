package config

type redisKey struct {
	UserInfo               string
	AccessToken            string
	YZM                    string
	UniDian                string
	Limit1S                string
	InitTopicUserFlowLimit string
}

var RedisKey = redisKey{
	UserInfo:               "mp2c:userinfo:%d",               // 变量1:用户id 使用 fmt.Sprintf(RedisKey.UserInfo,"1")
	AccessToken:            "mp2c:access_token:%s:%s",        // 变量1:应用平台 变量2:应用id
	YZM:                    "yzm:",                           // 拼接用户id
	UniDian:                "unidian:",                       // 拼接手机号
	Limit1S:                "Limit1S:",                       // 拼接行数名称
	InitTopicUserFlowLimit: "mp2c:inittopicuserflowlimit:%d", //拼接用户id
}
