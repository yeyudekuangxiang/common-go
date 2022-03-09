package activity

import "mio/model"

type BocShareBonusRecord struct {
	Id        int64      `json:"id"`
	UserId    int64      `json:"userId"`
	Value     int64      `json:"value"` //金额 分
	CreatedAt model.Time `json:"createdAt"`
	UpdatedAt model.Time `json:"updatedAt"`
}

func (BocShareBonusRecord) TableName() string {
	return "boc_share_bonus_record"
}
