package entity

import (
	"mio/config"
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
	UserPositionOrdinary UserPosition = "ordinary" //普通用户
	UserPositionBlue     UserPosition = "blue"     //蓝v
	UserPositionYellow   UserPosition = "yellow"   //黄v
)

var IconMap = map[UserPosition]string{
	UserPositionOrdinary: "",
	UserPositionBlue:     config.Config.OSS.CdnDomain + "/static/mp2c/images/user/blue.png",
	UserPositionYellow:   config.Config.OSS.CdnDomain + "/static/mp2c/user/positionIcon/oy_BA5ESvDiKAaKJ4GdrCOxkqP_4/d7e78457-5136-48ae-b64e-4c260e2a0c3a.png",
}

// https://resources.miotech.com/static/mp2c/event/cert/oy_BA5ESvDiKAaKJ4GdrCOxkqP_4/275cd192-6d9f-46e1-8023-ca8b16bd48fa.png

type UserGender string

const (
	UserGenderMale   UserGender = "MALE"
	UserGenderFemale UserGender = "FEMALE"
)

type Partner int

const (
	PartnerLoHoJa Partner = 1 //乐活家
)

var IconPartnerMap = map[Partner]string{
	PartnerLoHoJa: config.Config.OSS.CdnDomain + "/static/mp2c/event/cert/oy_BA5ESvDiKAaKJ4GdrCOxkqP_4/275cd192-6d9f-46e1-8023-ca8b16bd48fa.png",
}

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
	Partners     Partner      `gorm:"partners" json:"partners"` //乐活家 1是 2否
	Position     UserPosition `json:"position"`                 //身份 blue:蓝v yellow:黄v ordinary:普通用户
	PositionIcon string       `json:"positionIcon"`
	Risk         int          `json:"risk"`
	ChannelId    int64        `gorm:"column:channel_id" json:"channel_id"`
	Ip           string       `json:"ip"`
	CityCode     string       `json:"city_code"`
	Status       int          `json:"status,omitempty"` //0全部 1正常 2禁言 3封号 //暂时不用
	Auth         int          `json:"auth,omitempty"`   //发帖权限 0无权限 1发帖+评论 2评论权限
	Introduction string       `json:"introduction"`
}

func (u User) ShortUser() ShortUser {
	return ShortUser{
		ID:           u.ID,
		OpenId:       u.OpenId,
		AvatarUrl:    u.AvatarUrl,
		Gender:       u.Gender,
		Nickname:     u.Nickname,
		Partners:     u.Partners,
		Position:     u.Position,
		PositionIcon: u.PositionIcon,
		Auth:         u.Auth,
		Introduction: u.Introduction,
		Time:         u.Time,
	}
}

type ShortUser struct {
	ID           int64        `gorm:"primary_key;column:id" json:"id"`
	OpenId       string       `gorm:"column:openid" json:"openId"`
	AvatarUrl    string       `gorm:"column:avatar_url" json:"avatarUrl"`
	Gender       UserGender   `gorm:"column:gender" json:"gender"`
	Nickname     string       `gorm:"column:nick_name" json:"nickname"`
	Partners     Partner      `gorm:"partners" json:"partners"`
	Position     UserPosition `json:"position"`
	PositionIcon string       `json:"positionIcon"`
	Introduction string       `json:"introduction"`
	Auth         int          `json:"auth"` //发帖权限 0无权限 1有权限
	Time         model.Time   `json:"time"`
}

func (ShortUser) TableName() string {
	return "user"
}
