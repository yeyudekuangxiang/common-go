package qnr

import (
	"mio/internal/pkg/core/context"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	"mio/internal/pkg/repository/repotypes"
)

func NewSubjectRepository(ctx *context.MioContext) *SubjectRepository {
	return &SubjectRepository{
		ctx: ctx,
	}
}

type SubjectRepository struct {
	ctx *context.MioContext
}

func (repo SubjectRepository) Save(transaction *qnrEntity.Subject) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo SubjectRepository) List(do repotypes.GetQuestSubjectGetListBy) ([]qnrEntity.Subject, error) {
	db := repo.ctx.DB.Model(qnrEntity.Subject{})
	if do.QnrId != 0 {
		db.Where("qnr_id = ?", do.QnrId)
	}
	db.Order("sort desc,id asc")
	list := make([]qnrEntity.Subject, 0)
	return list, db.Find(&list).Error
}
