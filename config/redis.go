package config

type redisKey struct {
	UserInfo string
}

var RedisKey = redisKey{
	UserInfo: "mp2c:userinfo:%s", // 变量1:用户id 使用 fmt.Sprintf(RedisKey.UserInfo,"1")
}
