package repotypes

import "mio/internal/pkg/model/entity"

type GetQuestAnswerGetListBy struct {
	Openid      string
	OrderSource entity.OrderSource
	Limit       int
	Offset      int
	QnrId       int64
}

type GetQuestOptionGetListBy struct {
	SubjectIds []int64
}

type GetQuestSubjectGetListBy struct {
	QnrId int64
}
