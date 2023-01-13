package community

import "mio/internal/pkg/model/entity"

const (
	TopicTypeAll      = -1
	TopicTypeArticle  = 0
	TopicTypeActivity = 1
)

const (
	SignupStatusTrue  = 1
	SignupStatusFalse = 2
)

type FindAllActivitiesTagParams struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Offset      int                `json:"offset"`
	Limit       int                `json:"limit"` //limit为0时不限制数量
	OrderBy     entity.OrderByList `json:"orderBy"`
}

type FindAllActivitiesSignupParams struct {
	TopicId  int64  `json:"topicId"`
	UserId   int64  `json:"userId"`
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
	Gender   int    `json:"gender"`
	Age      int    `json:"age"`
	City     string `json:"city"`
	Wechat   string `json:"wechat"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type FindOneActivitiesSignupParams struct {
	Id           int64 `json:"id"`
	TopicId      int64 `json:"topicId"`
	UserId       int64 `json:"userId"`
	SignupStatus int   `json:"signupStatus"`
}
