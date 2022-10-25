package entity

import (
	"time"
)

type QuizDailyResult struct {
	ID           int64     `json:"id"`
	OpenId       string    `json:"openId" gorm:"column:openid"`
	AnswerDate   time.Time `json:"answerDate"`
	AnswerTime   time.Time `json:"answerTime"`
	CorrectNum   int       `json:"correctNum"`
	IncorrectNum int       `json:"incorrectNum"`
}
