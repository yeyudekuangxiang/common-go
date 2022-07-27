package entity

import "mio/internal/pkg/model"

type QuizSummary struct {
	ID               int64      `json:"id"`
	OpenId           string     `json:"openId" gorm:"column:openid"`
	TotalCorrectNum  int        `json:"totalCorrectNum"`
	TotalAnsweredNum int        `json:"totalAnsweredNum"`
	LastUpdateDate   model.Date `json:"LastUpdateDate"`
}
