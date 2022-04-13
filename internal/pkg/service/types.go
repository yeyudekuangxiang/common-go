package service

import (
	"mio/internal/pkg/model"
	entity2 "mio/internal/pkg/model/entity"
)

type TopicDetail struct {
	entity2.Topic
	IsLike        bool   `json:"isLike"`
	UpdatedAtDate string `json:"updatedAtDate"` //03-01
}

type CreatePointTransactionParam struct {
	OpenId       string
	Type         entity2.PointTransactionType
	Value        int
	AdditionInfo string
}
type CreateUserParam struct {
	OpenId      string             `json:"openId"`
	AvatarUrl   string             `json:"avatarUrl"`
	Gender      string             `json:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Nickname    string             `json:"nickname"`
	Birthday    model.Date         `json:"birthday"`
	PhoneNumber string             `json:"phoneNumber"`
	Source      entity2.UserSource `json:"source" binding:"oneof=mio mobile"`
	UnionId     string             `json:"unionId"`
}

type unidianTypeId struct {
	Test     string
	FiveYuan string
}

var UnidianTypeId = unidianTypeId{
	Test:     "10013", // 测试
	FiveYuan: "10689", // 5元话费
}

type SubmitOrderParam struct {
	Order SubmitOrder
	Items []SubmitOrderItem
}
type SubmitOrder struct {
	AddressId string
	UserId    int64
	OrderType entity2.OrderType
}
type SubmitOrderItem struct {
	ItemId string
	Count  int
}

type submitOrderParam struct {
	Order submitOrder
	Items []submitOrderItem
}
type submitOrder struct {
	AddressId string
	UserId    int64
	TotalCost int
	OrderType entity2.OrderType
}
type submitOrderItem struct {
	ItemId string
	Count  int
	Cost   int
}
type SubmitOrderForGreenParam struct {
	AddressId string
	UserId    int64
	ItemId    string
}
type SubmitOrderForActivityParam struct {
	AddressId string
	UserId    int64
	ItemId    string
	Activity  string
}
type CalculateProductResult struct {
	TotalCost int
	ItemList  []submitOrderItem
}
type ExchangeCallbackResult struct {
	BizId   string
	Credits int
}

type AutoLoginParam struct {
	UserId   int64
	Path     string
	DCustom  string
	Transfer string
	SignKeys string
}
type AutoLoginOpenIdParam struct {
	UserId   int64
	OpenId   string
	Path     string
	DCustom  string
	Transfer string
	SignKeys string
}
type BindPhoneByIVParam struct {
	UserId        int64
	EncryptedData string
	IV            string
}
