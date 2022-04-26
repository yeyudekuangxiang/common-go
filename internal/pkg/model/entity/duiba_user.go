package entity

import "mio/internal/pkg/model"

type DuiBaUser struct {
	ID         int64      `json:"id"`
	OpenId     string     `json:"openId" gorm:"column:openid"`
	ActivityId string     `json:"activityId"`
	IsNew      int        `json:"isNew"`
	CreateTime model.Time `json:"createTime"`
	CreatedAt  model.Time `json:"createdAt"`
	UpdatedAt  model.Time `json:"updatedAt"`
}

func (DuiBaUser) TableName() string {
	return "duiba_user"
}
