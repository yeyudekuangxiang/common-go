package qnr

import "mio/internal/pkg/model"

type Option struct {
	ID             int64
	Title          string
	Sort           int8
	SubjectId      model.LongID
	Remind         string
	JumpSubject    int64
	RelatedSubject string
}

func (Option) TableName() string {
	return "qnr_option"
}
