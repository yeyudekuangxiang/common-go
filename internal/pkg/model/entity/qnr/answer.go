package qnr

import "mio/internal/pkg/model"

type Answer struct {
	ID        int64
	QnrId     int64
	SubjectId model.LongID
	UserId    int64
	Answer    string
}

func (Answer) TableName() string {
	return "qnr_answer"
}
