package repotypes

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
)

type GetQuestAnswerGetListBy struct {
	Openid      string
	OrderSource entity.OrderSource
	Limit       int
	Offset      int
	QnrId       int64
}

type GetQuestOptionGetListBy struct {
	SubjectIds []model.LongID
}

type GetQuestOptionGetListByUid struct {
	Uid   int64
	QnrId int64
}

type GetQuestSubjectGetListBy struct {
	QnrId int64
}

type GetQuestUserGetById struct {
	UserId int64
	OpenId string
}
