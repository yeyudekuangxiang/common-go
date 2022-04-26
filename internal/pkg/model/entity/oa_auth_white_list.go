package entity

import "mio/internal/pkg/model"

type OaAuthWhite struct {
	ID        int `json:"id"`
	Domain    string
	AppId     string `gorm:"column:appid"`
	CreatedAt model.Time
	UpdatedAt model.Time
}

func (OaAuthWhite) TableName() string {
	return "oa_auth_white_list"
}
