package model

import (
	"encoding/json"
)

type OrderStatus string

const (
	OrderStatusWaitAudit OrderStatus = "waitAudit"
	OrderStatusWaitSend  OrderStatus = "waitSend"
	OrderStatusAfterSend OrderStatus = "afterSend"
	OrderStatusSuccess               = "success"
	OrderStatusFail                  = "fail"
)

type OrderType string

const (
	OrderTypeObject    OrderType = "object"
	OrderTypeCoupon    OrderType = "coupon"
	OrderTypeVirtual   OrderType = "virtual"
	OrderTypeAlipay    OrderType = "alipay"
	OrderTypePhoneBill OrderType = "phonebill"
	OrderTypePhoneQB   OrderType = "qb"
)

type OrderInfo struct {
	Base
	OrderNum         string             `json:"orderNum" form:"orderNum" binding:"required" alias:"orderNum"`
	DevelopBizId     string             `json:"developBizId" form:"developBizId"`
	CreateTime       string             `json:"createTime" form:"createTime" binding:"required" alias:"createTime"`
	FinishTime       string             `json:"finishTime" form:"finishTime" binding:"required" alias:"finishTime"`
	TotalCredits     string             `json:"totalCredits" form:"totalCredits"  alias:"totalCredits"`
	ConsumerPayPrice string             `json:"consumerPayPrice" form:"consumerPayPrice" binding:"required" alias:"consumerPayPrice"`
	Source           string             `json:"source" form:"source" binding:"required" alias:"source"`
	OrderStatus      OrderStatus        `json:"orderStatus" form:"orderStatus" binding:"required" alias:"orderStatus"`
	ErrorMsg         string             `json:"errorMsg" form:"errorMsg"`
	Type             OrderType          `json:"type" form:"type" binding:"required" alias:"type"`
	ExpressPrice     string             `json:"expressPrice" form:"expressPrice"`
	Account          string             `json:"account" form:"account" binding:"required" alias:"account"`
	OrderItemList    OrderItemListStr   `json:"orderItemList" form:"orderItemList" binding:"required" alias:"orderItemList"`
	ReceiveAddrInfo  ReceiveAddrInfoStr `json:"receiveAddrInfo" form:"receiveAddrInfo"`
}

func (o OrderInfo) ToMap() map[string]string {
	return map[string]string{
		"uid":              o.Uid,
		"appKey":           o.AppKey,
		"timestamp":        o.Timestamp,
		"sign":             o.Sign,
		"orderNum":         o.OrderNum,
		"developBizId":     o.DevelopBizId,
		"createTime":       o.CreateTime,
		"finishTime":       o.FinishTime,
		"totalCredits":     o.TotalCredits,
		"consumerPayPrice": o.ConsumerPayPrice,
		"source":           o.Source,
		"orderStatus":      string(o.OrderStatus),
		"errorMsg":         o.ErrorMsg,
		"type":             string(o.Type),
		"expressPrice":     o.ExpressPrice,
		"account":          o.Account,
		"orderItemList":    string(o.OrderItemList),
		"receiveAddrInfo":  string(o.ReceiveAddrInfo),
	}
}

type OrderItemListStr string

func (o OrderItemListStr) OrderItemList() ([]OrderItem, error) {
	list := make([]OrderItem, 0)
	if o == "" {
		return list, nil
	}

	err := json.Unmarshal([]byte(o), &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

type ReceiveAddrInfoStr string

func (r ReceiveAddrInfoStr) ReceiveAddrInfo() (*OrderAddressInfo, error) {
	address := OrderAddressInfo{}
	if r == "" {
		return &address, nil
	}

	return &address, json.Unmarshal([]byte(r), &address)
}
