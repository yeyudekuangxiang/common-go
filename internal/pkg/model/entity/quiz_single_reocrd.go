package entity

import (
	"time"
)

type QuizSingleRecord struct {
	ID         int64     `json:"id"`
	OpenId     string    `json:"openId" gorm:"column:openid"`
	QuestionId string    `json:"questionId"`
	Correct    bool      `json:"correct"`
	AnswerTime time.Time `json:"answerTime"`
	AnswerDate time.Time `json:"answerDate"`
}
