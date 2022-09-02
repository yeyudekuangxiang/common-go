package qnr

import (
	"mio/internal/pkg/core/context"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	"mio/internal/pkg/repository/repotypes"
)

func NewAnswerRepository(ctx *context.MioContext) *AnswerRepository {
	return &AnswerRepository{
		ctx: ctx,
	}
}

type AnswerRepository struct {
	ctx *context.MioContext
}

func (repo AnswerRepository) Save(transaction *qnrEntity.Answer) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo AnswerRepository) CreateInBatches(transaction []qnrEntity.Answer) error {
	return repo.ctx.DB.CreateInBatches(transaction, len(transaction)).Error
}

func (repo AnswerRepository) GetListBy(by repotypes.GetQuestOptionGetListBy) ([]qnrEntity.Answer, error) {
	list := make([]qnrEntity.Answer, 0)
	db := repo.ctx.DB.Model(qnrEntity.Answer{})
	if len(by.SubjectIds) != 0 {
		db.Where("subject_id in (?) ", by.SubjectIds)
	}
	db.Order("sort desc")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}

func (repo AnswerRepository) GetListByUid(by repotypes.GetQuestOptionGetListByUid) ([]qnrEntity.Answer, error) {
	list := make([]qnrEntity.Answer, 0)
	db := repo.ctx.DB.Model(qnrEntity.Answer{})
	if by.QnrId != 0 {
		db.Where("qnr_id", by.QnrId)
	}
	if by.Uid != 0 {
		db.Where("user_id", by.Uid)
	}
	db.Order("sort desc")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}
