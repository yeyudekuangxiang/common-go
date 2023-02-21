package entity

import duibaApi "gitlab.miotech.com/miotech-application/backend/common-go/duiba/api/model"

type DuiBaOrder struct {
	ID               int64
	OrderNum         string
	DevelopBizId     string
	CreateTime       int64
	FinishTime       int64
	TotalCredits     int
	ConsumerPayPrice float64
	Source           string
	OrderStatus      duibaApi.OrderStatus
	ErrorMsg         string
	Type             duibaApi.OrderType
	ExpressPrice     string
	Account          string
	OrderItemList    string
	ReceiveAddrInfo  string
	OrderId          string
	UserId           int64
}

func (DuiBaOrder) TableName() string {
	return "duiba_order"
}
