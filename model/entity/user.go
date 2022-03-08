package entity

type User struct {
	ID          int64  `gorm:"primary_key;column:id" json:"id"`
	OpenId      string `gorm:"column:openid" json:"openId"`
	AvatarUrl   string `gorm:"column:avatar_url" json:"avatarUrl"`
	Gender      string `gorm:"column:gender" json:"gender"`
	Nickname    string `gorm:"column:nick_name" json:"nickname"`
	Birthday    string `gorm:"column:birthday" json:"birthday"`
	PhoneNumber string `gorm:"column:phone_number" json:"phoneNumber"`
}
type ShortUser struct {
	ID        int64  `gorm:"primary_key;column:id" json:"id"`
	OpenId    string `gorm:"column:openid" json:"openId"`
	AvatarUrl string `gorm:"column:avatar_url" json:"avatarUrl"`
	Nickname  string `gorm:"column:nick_name" json:"nickname"`
}
