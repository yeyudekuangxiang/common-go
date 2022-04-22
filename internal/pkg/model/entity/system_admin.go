package entity

import "mio/internal/pkg/model"

type SystemAdmin struct {
	ID        int        `gorm:"primary_key;column:id" json:"id"`
	Nickname  string     `json:"nickname"`
	RealName  string     `json:"realName" gorm:"column:realname"`
	Avatar    string     `json:"avatar"`
	Status    int        `json:"status"` //1正常 2已停用(已离职)
	Account   string     `json:"account"`
	Password  string     `json:"-"`
	Phone     string     `json:"phone"`
	CreatedAt model.Time `json:"createdAt"`
	UpdatedAt model.Time `json:"updatedAt"`
}
