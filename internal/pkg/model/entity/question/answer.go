package question

import "mio/internal/pkg/model"

type Answer struct {
	ID        int64
	QnrId     int64
	SubjectId model.LongID
	UserId    model.LongID
	Answer    string
	Carbon    float64
}

func (Answer) TableName() string {
	return "question_answer"
}
