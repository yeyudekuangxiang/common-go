package entity

import "mio/internal/pkg/model"

type QuizDailyResult struct {
	ID           int64      `json:"id"`
	OpenId       string     `json:"openId" gorm:"column:openid"`
	AnswerDate   model.Date `json:"answerDate"`
	AnswerTime   model.Time `json:"answerTime"`
	CorrectNum   int        `json:"correctNum"`
	IncorrectNum int        `json:"incorrectNum"`
}
