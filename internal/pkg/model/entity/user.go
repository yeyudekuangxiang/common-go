package entity

import (
	"mio/internal/pkg/model"
)

type UserSource string

const (
	UserSourceMio      UserSource = "mio"
	UserSourceMobile   UserSource = "mobile"
	UserSourceMioSrvOA UserSource = "mio-srv-oa"
	UserSourceMioSubOA UserSource = "mio-sub-oa"
)

type UserPosition string

const (
	UserPositionBlue   UserPosition = "blue"   //蓝v
	UserPositionYellow UserPosition = "yellow" //黄v
	UserPositionLohoja UserPosition = "lohoja" //乐活家
)

type UserGender string

const (
	UserGenderMale   UserGender = "MALE"
	UserGenderFemale UserGender = "FEMALE"
)

type User struct {
	ID           int64        `gorm:"primary_key;column:id" json:"id"`
	OpenId       string       `gorm:"column:openid" json:"openId"`
	AvatarUrl    string       `gorm:"column:avatar_url" json:"avatarUrl"`
	Gender       UserGender   `gorm:"column:gender" json:"gender"`
	Nickname     string       `gorm:"column:nick_name" json:"nickname"`
	Birthday     model.Date   `gorm:"column:birthday" json:"birthday"`
	PhoneNumber  string       `gorm:"column:phone_number" json:"phoneNumber"`
	Source       UserSource   `gorm:"column:source" json:"source"` //用户来源 mio(绿喵小程序) mobile(手机号注册)
	UnionId      string       `gorm:"column:unionid" json:"unionId"`
	Time         model.Time   `gorm:"time" json:"time"`
	GUID         string       `gorm:"guid" json:"guid"`
	Position     UserPosition `json:"position"`
	PositionIcon string       `json:"positionIcon"`
	Risk         int          `json:"risk"`
	ChannelId    int64        `gorm:"column:channel_id" json:"channel_id"`
	State        int          `json:"state" gorm:"default:0"` //用户状态 0正常 1禁言 2封号
}

func (u User) ShortUser() ShortUser {
	return ShortUser{
		ID:           u.ID,
		AvatarUrl:    u.AvatarUrl,
		Gender:       u.Gender,
		Nickname:     u.Nickname,
		Position:     u.Position,
		PositionIcon: u.PositionIcon,
	}
}

type ShortUser struct {
	ID           int64        `gorm:"primary_key;column:id" json:"id"`
	AvatarUrl    string       `gorm:"column:avatar_url" json:"avatarUrl"`
	Gender       UserGender   `gorm:"column:gender" json:"gender"`
	Nickname     string       `gorm:"column:nick_name" json:"nickname"`
	Position     UserPosition `json:"position"`
	PositionIcon string       `json:"positionIcon"`
}

func (ShortUser) TableName() string {
	return "user"
}
