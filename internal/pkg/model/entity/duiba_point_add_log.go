package entity

import (
	"mio/internal/pkg/model"
	duibaApi "mio/pkg/duiba/api/model"
)

type DuiBaPointAddLog struct {
	ID            int64
	Uid           string
	Credits       int64
	Type          duibaApi.PointAddType
	OrderNum      string
	SubOrderNum   string
	Timestamp     int64
	Description   string
	Ip            string
	Sign          string
	AppKey        string
	CreatedAt     model.Time
	TransactionId string
}

func (DuiBaPointAddLog) TableName() string {
	return "duiba_point_add_log"
}
