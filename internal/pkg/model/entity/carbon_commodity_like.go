package entity

import "time"

type CarbonCommodityLike struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CommodityId int64     `gorm:"type:int8;not null" json:"commodityId"`
	Status      int       `gorm:"type:int2;not null" json:"status"`
	UserId      int64     `gorm:"type:int8;not null" json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (c CarbonCommodityLike) TableName() string {
	return "carbon_commodity_like"
}
