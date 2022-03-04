package entity

type User struct {
	ID          int    `gorm:"primary_key;column:id" json:"id"`
	OpenId      string `gorm:"column:openid"`
	AvatarUrl   string `gorm:"column:avatar_url"`
	Gender      string `gorm:"column:gender"`
	Nickname    string `gorm:"column:nick_name"`
	Birthday    string `gorm:"column:birthday"`
	PhoneNumber string `gorm:"column:phone_number"`
}
type ShortUser struct {
	ID        int    `gorm:"primary_key;column:id" json:"id"`
	OpenId    string `gorm:"column:openid"`
	AvatarUrl string `gorm:"column:avatar_url"`
	Nickname  string `gorm:"column:nick_name"`
}
