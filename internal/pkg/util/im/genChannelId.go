package im

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"strconv"
)

func GenChannelIdByFriend(userId, friendId int64) string {
	num := userId + friendId
	return genChannelId(strconv.FormatInt(num, 10))
}

func GenChannelIdByGroup(groupName string) string {
	return encrypttool.Md5(groupName)
}

func genChannelId(key string) string {
	return encrypttool.Md5(key)
}
