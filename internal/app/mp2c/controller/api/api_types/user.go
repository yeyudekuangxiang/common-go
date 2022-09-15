package api_types

import (
	"mio/internal/pkg/model"
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

type UserInfoVO struct {
	ID           int64               `gorm:"primary_key;column:id" json:"id"`
	OpenId       string              `gorm:"column:openid" json:"openId"`
	AvatarUrl    string              `gorm:"column:avatar_url" json:"avatarUrl"`
	Gender       entity.UserGender   `gorm:"column:gender" json:"gender"`
	Nickname     string              `gorm:"column:nick_name" json:"nickname"`
	Birthday     model.Date          `gorm:"column:birthday" json:"birthday"`
	PhoneNumber  string              `gorm:"column:phone_number" json:"phoneNumber"`
	Source       entity.UserSource   `gorm:"column:source" json:"source"` //用户来源 mio(绿喵小程序) mobile(手机号注册)
	UnionId      string              `gorm:"column:unionid" json:"unionId"`
	Time         model.Time          `gorm:"time" json:"time"`
	GUID         string              `gorm:"guid" json:"guid"`
	Partners     entity.Partner      `gorm:"partners" json:"partners"` //合作商 0:全部 1:乐活家 2:非乐活家
	Position     entity.UserPosition `json:"position"`                 //身份 blue:蓝v yellow:黄v ordinary:普通用户
	PositionIcon string              `json:"positionIcon"`
	Risk         int                 `json:"risk"`
	ChannelId    int64               `gorm:"column:channel_id" json:"channel_id"`
	Ip           string              `json:"ip"`
	CityCode     string              `json:"city_code"`
	Status       int                 `json:"status"` //0全部 1正常 2禁言 3封号 //暂时不用
	Auth         int                 `json:"auth"`   //发帖权限 0无权限 1有权限
	ChannelName  string              `json:"channelName"`
}
