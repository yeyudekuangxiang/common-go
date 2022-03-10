package service

import (
	"mio/model"
	"mio/model/entity"
)

type TopicDetail struct {
	entity.Topic
	IsLike bool `json:"isLike"`
}

type CreatePointTransactionParam struct {
	OpenId       string
	Type         entity.PointTransactionType
	Value        int
	AdditionInfo string
}
type CreateUserParam struct {
	OpenId      string            `json:"openId"`
	AvatarUrl   string            `json:"avatarUrl"`
	Gender      string            `json:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Nickname    string            `json:"nickname"`
	Birthday    model.Date        `json:"birthday"`
	PhoneNumber string            `json:"phoneNumber"`
	Source      entity.UserSource `json:"source" binding:"oneof=mio mobile"`
	UnionId     string            `json:"unionId"`
}

type unidianTypeId struct {
	Test     string
	FiveYuan string
}

var UnidianTypeId = unidianTypeId{
	Test:     "10013", // 测试
	FiveYuan: "10689", // 5元话费
}
