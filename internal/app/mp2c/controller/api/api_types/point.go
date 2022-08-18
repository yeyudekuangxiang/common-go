package api_types

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
)

type PointRecordInfo struct {
	ID             int64                       `json:"id"`
	TransactionId  string                      `json:"transactionId"`
	Type           entity.PointTransactionType `json:"type"`
	TypeText       string                      `json:"typeText"`
	Value          int64                       `json:"value"`
	CreateTime     model.Time                  `json:"createTime"`
	TimeStr        string                      `json:"timeStr"`
	AdditionalInfo entity.AdditionalInfo       `json:"additionalInfo"`
	AdminId        int                         `json:"adminId"`
	Note           string                      `json:"note"`
}
