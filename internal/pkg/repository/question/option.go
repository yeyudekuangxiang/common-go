package question

import (
	"mio/internal/pkg/core/context"
	questionEntity "mio/internal/pkg/model/entity/question"
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

func (repo OptionRepository) Save(transaction *questionEntity.Option) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo OptionRepository) GetListBy(by repotypes.GetQuestOptionGetListBy) ([]questionEntity.Option, error) {
	list := make([]questionEntity.Option, 0)
	db := repo.ctx.DB.Model(questionEntity.Option{})
	if len(by.SubjectIds) != 0 {
		db.Where("subject_id in (?) ", by.SubjectIds)
	}
	db.Order("sort desc,id asc")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}

func (repo OptionRepository) CreateInBatches(transaction []questionEntity.Option) error {
	return repo.ctx.DB.CreateInBatches(transaction, len(transaction)).Error
}
