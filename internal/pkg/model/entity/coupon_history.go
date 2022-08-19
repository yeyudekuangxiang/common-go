package entity

import "mio/internal/pkg/model"

type CouponHistory struct {
	ID           int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	OpenId       string
	CouponTypeId int64
	CouponType   string
	Code         string
	CreateTime   model.CreatedTime
	UpdateTime   model.UpdatedTime
}
