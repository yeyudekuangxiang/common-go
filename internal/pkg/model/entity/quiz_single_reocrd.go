package entity

import "mio/internal/pkg/model"

type QuizSingleRecord struct {
	ID         int64      `json:"id"`
	OpenId     string     `json:"openId" gorm:"column:openid"`
	QuestionId string     `json:"questionId"`
	Correct    bool       `json:"correct"`
	AnswerTime model.Time `json:"answerTime"`
	AnswerDate model.Date `json:"answerDate"`
}
