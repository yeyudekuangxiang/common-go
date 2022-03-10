package config

type redisKey struct {
	UserInfo string
	YZM      string
	UniDian  string
	Limit1S  string
}

var RedisKey = redisKey{
	UserInfo: "mp2c:userinfo:%s", // 变量1:用户id 使用 fmt.Sprintf(RedisKey.UserInfo,"1")
	YZM:      "yzm:",             // 拼接用户id
	UniDian:  "unidian:",         // 拼接手机号
	Limit1S:  "Limit1S:",         // 拼接行数名称
}
