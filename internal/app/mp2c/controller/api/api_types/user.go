package api_types

import (
	"mio/internal/pkg/model/entity"
)

type UserVO struct {
	ID          int64             ` json:"id"`
	OpenId      string            `json:"openId"`
	AvatarUrl   string            `json:"avatarUrl"`
	Nickname    string            `json:"nickname"`
	PhoneNumber string            `json:"phoneNumber"`
	Source      entity.UserSource `json:"source"` //用户来源 mio(绿喵小程序) mobile(手机号注册)
	RegTime     string            `json:"regTime"`
	Risk        int               `json:"risk"`
	ChannelName string            `json:"channelName"`
	CityName    string            `json:"cityCode"`
	Status      int               `json:"status"` //0全部 1正常 2禁言 3封号 //暂时不用
	Point       int64             `json:"point"`
}
