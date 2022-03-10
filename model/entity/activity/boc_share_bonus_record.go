package activity

import "mio/model"

type BocShareBonusType string

const (
	BocShareBonusMio   = "mio"
	BocShareBonusBoc   = "boc"
	BocShareBonusBoc10 = "boc10"
)

type BocShareBonusRecord struct {
	Id        int64             `json:"id"`
	UserId    int64             `json:"userId"`
	Value     int64             `json:"value"` //金额 分
	Type      BocShareBonusType `json:"type"`
	Info      string            `json:"info"` //描述信息
	CreatedAt model.Time        `json:"createdAt"`
	UpdatedAt model.Time        `json:"updatedAt"`
}

func (BocShareBonusRecord) TableName() string {
	return "boc_share_bonus_record"
}
