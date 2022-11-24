package im

import (
	"mio/internal/pkg/util/encrypt"
	"strconv"
)

func GenChannelIdByFriend(userId, friendId int64) string {
	num := userId + friendId
	return genChannelId(strconv.FormatInt(num, 10))
}

func GenChannelIdByGroup(groupName string) string {
	return encrypt.Md5(groupName)
}

func genChannelId(key string) string {
	return encrypt.Md5(key)
}
