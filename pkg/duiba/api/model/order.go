package model

import (
	"encoding/json"
	"log"
	"strconv"
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
	CreateTime       IntStr             `json:"createTime" form:"createTime" binding:"required" alias:"createTime"`
	FinishTime       IntStr             `json:"finishTime" form:"finishTime" binding:"required" alias:"finishTime"`
	TotalCredits     IntStr             `json:"totalCredits" form:"totalCredits"  alias:"totalCredits"`
	ConsumerPayPrice FloatStr           `json:"consumerPayPrice" form:"consumerPayPrice" binding:"required" alias:"consumerPayPrice"`
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
		"timestamp":        string(o.Timestamp),
		"sign":             o.Sign,
		"orderNum":         o.OrderNum,
		"developBizId":     o.DevelopBizId,
		"createTime":       string(o.CreateTime),
		"finishTime":       string(o.FinishTime),
		"totalCredits":     string(o.TotalCredits),
		"consumerPayPrice": string(o.ConsumerPayPrice),
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

func (o OrderItemListStr) OrderItemList() []OrderItem {
	list := make([]OrderItem, 0)
	if o == "" {
		return list
	}

	err := json.Unmarshal([]byte(o), &list)
	log.Println(err)
	return list
}
func (o OrderItemListStr) OrderItemListE() ([]OrderItem, error) {
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

type FloatStr string

func (f FloatStr) ToFloat() float64 {
	data, _ := strconv.ParseFloat(string(f), 64)
	return data
}
func (f FloatStr) ToFloatE() (float64, error) {
	return strconv.ParseFloat(string(f), 64)
}

type IntStr string

func (i IntStr) ToInt() int64 {
	data, _ := strconv.ParseInt(string(i), 10, 64)
	return data
}
func (i IntStr) ToIntE() (int64, error) {
	return strconv.ParseInt(string(i), 10, 64)
}
