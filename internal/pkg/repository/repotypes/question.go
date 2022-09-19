package repotypes

import (
	"mio/internal/pkg/model"
)

type GetQuestionOptionGetListBy struct {
	SubjectIds []model.LongID
}

type GetQuestionOptionGetListByUid struct {
	Uid        int64
	QuestionId int64
}

type GetQuestionSubjectGetListBy struct {
	QuestionId int64
}

type GetQuestionUserGetById struct {
	UserId int64
	OpenId string
}

type GetQuestionUserCarbon struct {
	Uid        int64
	QuestionId int64
}

type UserAnswerStruct struct {
	CategoryId string  `json:"category_id"`
	Carbon     float64 `json:"carbon"`
}
