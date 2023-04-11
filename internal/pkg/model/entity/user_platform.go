package entity

import "time"

type UserPlatform struct {
	Id        int64     `gorm:"primaryKey;autoIncrement;column:id"`
	Guid      string    `gorm:"column:guid"`       // guid
	Openid    string    `gorm:"column:openid"`     // openid
	Nickname  string    `gorm:"column:nickname"`   // 昵称
	AvatarUrl string    `gorm:"column:avatar_url"` // 头像
	Sex       int64     `gorm:"column:sex"`        // 性别 0未知 1男性 2女性
	Unionid   string    `gorm:"column:unionid"`    // unionid
	Country   string    `gorm:"column:country"`    // 国家
	Province  string    `gorm:"column:province"`   // 省份
	City      string    `gorm:"column:city"`       // 国家
	Ip        string    `gorm:"column:ip"`         // ip地址
	Platform  string    `gorm:"column:platform"`   // wechat(微信app)、wxminiapp(微信小程序)、wxoa(公众号)
	CreatedAt time.Time `gorm:"column:created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at"` // 更新时间
}

func (UserPlatform) TableName() string {
	return "user_platform"
}

type UserPlatformType string

const (
	// UserPlatformWechat 微信app
	UserPlatformWechat UserPlatformType = "wechat"
	// UserPlatformWxMiniApp 微信小程序
	UserPlatformWxMiniApp UserPlatformType = "wxminiapp"
	// UserPlatformWxOA 公众号
	UserPlatformWxOA UserPlatformType = "wxoa"
)
