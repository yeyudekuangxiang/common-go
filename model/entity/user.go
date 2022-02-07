package entity

type User struct {
	ID        int    `gorm:"primary_key;column:id" json:"id"`
	NickName  string `gorm:"column:nick_name" json:"nick_name"`
	AvatarUrl string `gorm:"avatar_url" json:"avatar_url"`
}
