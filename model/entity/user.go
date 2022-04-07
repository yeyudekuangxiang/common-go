package entity

import (
	"mio/model"
)

type UserSource string

const (
	UserSourceMio    UserSource = "mio"
	UserSourceMobile UserSource = "mobile"
)

type User struct {
	ID          int64      `gorm:"primary_key;column:id" json:"id"`
	OpenId      string     `gorm:"column:openid" json:"openId"`
	AvatarUrl   string     `gorm:"column:avatar_url" json:"avatarUrl"`
	Gender      string     `gorm:"column:gender" json:"gender"`
	Nickname    string     `gorm:"column:nick_name" json:"nickname"`
	Birthday    model.Date `gorm:"column:birthday" json:"birthday"`
	PhoneNumber string     `gorm:"column:phone_number" json:"phoneNumber"`
	Source      UserSource `gorm:"column:source" json:"source"` //用户来源 mio(绿喵小程序) mobile(手机号注册)
	UnionId     string     `gorm:"column:unionid" json:"unionId"`
	Time        model.Time `gorm:"time" json:"time"`
}

func (u User) ShortUser() ShortUser {
	return ShortUser{
		ID:        u.ID,
		AvatarUrl: u.AvatarUrl,
		Gender:    u.Gender,
		Nickname:  u.Nickname,
	}
}

type ShortUser struct {
	ID        int64  `gorm:"primary_key;column:id" json:"id"`
	AvatarUrl string `gorm:"column:avatar_url" json:"avatarUrl"`
	Gender    string `gorm:"column:gender" json:"gender"`
	Nickname  string `gorm:"column:nick_name" json:"nickname"`
}

func (ShortUser) TableName() string {
	return "user"
}
