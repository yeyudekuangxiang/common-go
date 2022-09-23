package question

import "mio/internal/pkg/model"

type Option struct {
	ID             int64
	Title          string
	Sort           int8
	Carbon         float64
	SubjectId      model.LongID
	Remind         string
	JumpSubject    int64
	RelatedSubject string
}

func (Option) TableName() string {
	return "question_option"
}
