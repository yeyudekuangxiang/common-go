package qnr

import (
	"mio/internal/pkg/core/context"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	"mio/internal/pkg/repository/repotypes"
)

func NewOptionRepository(ctx *context.MioContext) *OptionRepository {
	return &OptionRepository{
		ctx: ctx,
	}
}

type OptionRepository struct {
	ctx *context.MioContext
}

func (repo OptionRepository) Save(transaction *qnrEntity.Option) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo OptionRepository) GetListBy(by repotypes.GetQuestOptionGetListBy) ([]qnrEntity.Option, error) {
	list := make([]qnrEntity.Option, 0)
	db := repo.ctx.DB.Model(qnrEntity.Option{})
	if len(by.SubjectIds) != 0 {
		db.Where("subject_id in (?) ", by.SubjectIds)
	}
	db.Order("sort desc,id asc")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}
