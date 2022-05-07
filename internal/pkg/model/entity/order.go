package entity

import (
	"mio/internal/pkg/model"
)

type OrderStatus string

const (
	OrderStatusPaid      = "PAID"       //已支付
	OrderStatusPending   = "PENDING"    //
	OrderStatusInTransit = "IN_TRANSIT" //已发货
	OrderStatusComplete  = "COMPLETED"  //已完成
	OrderStatusError     = "ERROR"      //订单异常
)

type OrderType string

const (
	OrderTypePurchase   = "PURCHASE" //下单购买
	OrderTypeRedeem     = "REDEEM"   //兑换?
	OrderTypeGreenTorch = "GREEN_TORCH"
)

type OrderSource string

const (
	OrderSourceMio   OrderSource = "mio"
	OrderSourceDuiBa OrderSource = "duiba"
)

type Order struct {
	ID               int64       `gorm:"primary_key" json:"id"`
	OrderId          string      `json:"orderId"`
	AddressId        *string     `json:"addressId"`
	OpenId           string      `gorm:"column:openid" json:"openid"`
	TotalCost        int         `json:"totalCost"`
	Status           OrderStatus `json:"orderStatus"`
	PaidTime         model.Time  `json:"paidTime"`
	InTransitTime    model.Time  `json:"inTransitTime"`
	CompletedTime    model.Time  `json:"completedTime"`
	TrackingNumber   string      `json:"trackingNumber"`
	OrderReferenceId string      `json:"orderReferenceId"`
	OrderType        OrderType   `json:"orderType"`
	Source           OrderSource `json:"source"`
	ThirdOrderNo     string      `json:"thirdOrderNo"`
}

func (order Order) ShortOrder() ShortOrder {
	return ShortOrder{
		OrderId:          order.OrderId,
		TotalCost:        order.TotalCost,
		Status:           order.Status,
		PaidTime:         order.PaidTime,
		TrackingNumber:   order.TrackingNumber,
		OrderReferenceId: order.OrderReferenceId,
		OrderType:        order.OrderType,
	}
}

type ShortOrder struct {
	OrderId          string      `json:"orderId"`
	TotalCost        int         `json:"totalCost"`
	Status           OrderStatus `json:"orderStatus"`
	PaidTime         model.Time  `json:"paidTime"`
	TrackingNumber   string      `json:"trackingNumber"`
	OrderReferenceId string      `json:"orderReferenceId"`
	OrderType        OrderType   `json:"orderType"`
}
