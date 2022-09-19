package repotypes

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity/question"
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
	CategoryId question.QuestionCategoryType `json:"categoryId"`
	Carbon     float64                       `json:"carbon"`
}
