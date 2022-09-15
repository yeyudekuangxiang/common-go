package question

import (
	"mio/internal/pkg/core/context"
	questionEntity "mio/internal/pkg/model/entity/question"
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

func (repo SubjectRepository) Save(transaction *questionEntity.Subject) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo SubjectRepository) List(do repotypes.GetQuestSubjectGetListBy) ([]questionEntity.Subject, error) {
	db := repo.ctx.DB.Model(questionEntity.Subject{})
	if do.QnrId != 0 {
		db.Where("qnr_id = ?", do.QnrId)
	}
	db.Order("sort desc,id asc")
	list := make([]questionEntity.Subject, 0)
	return list, db.Find(&list).Error
}

func (repo SubjectRepository) CreateInBatches(transaction []questionEntity.Subject) error {
	return repo.ctx.DB.CreateInBatches(transaction, len(transaction)).Error
}
