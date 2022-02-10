package entity

type Admin struct {
	ID       int    `gorm:"primary_key;column:id" json:"id"`
	UName    string `gorm:"column:uname" json:"uName"`
	RealName string `gorm:"column:realname" json:"realName"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
}
