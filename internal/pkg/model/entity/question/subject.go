package question

import "mio/internal/pkg/model"

type Subject struct {
	ID         int64
	SubjectId  model.LongID
	QnrId      int64
	CategoryId int64
	Title      string
	Type       int8
	IsHide     int8
	Remind     string
	Sort       int8
}

func (Subject) TableName() string {
	return "question_subject"
}
