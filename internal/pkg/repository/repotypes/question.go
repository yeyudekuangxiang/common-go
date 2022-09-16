package repotypes

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
)

type GetQuestionAnswerGetListBy struct {
	Openid      string
	OrderSource entity.OrderSource
	Limit       int
	Offset      int
	QuestionId  int64
}

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
