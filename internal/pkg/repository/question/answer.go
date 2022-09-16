package question

import (
	"mio/internal/pkg/core/context"
	QuestionEntity "mio/internal/pkg/model/entity/question"
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

func (repo AnswerRepository) Save(transaction *QuestionEntity.Answer) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo AnswerRepository) CreateInBatches(transaction []QuestionEntity.Answer) error {
	return repo.ctx.DB.CreateInBatches(transaction, len(transaction)).Error
}

func (repo AnswerRepository) GetListBy(by repotypes.GetQuestOptionGetListBy) ([]QuestionEntity.Answer, error) {
	list := make([]QuestionEntity.Answer, 0)
	db := repo.ctx.DB.Model(QuestionEntity.Answer{})
	if len(by.SubjectIds) != 0 {
		db.Where("subject_id in (?) ", by.SubjectIds)
	}
	db.Order("sort desc")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}

func (repo AnswerRepository) GetListByUid(by repotypes.GetQuestionOptionGetListByUid) ([]QuestionEntity.Answer, error) {
	list := make([]QuestionEntity.Answer, 0)
	db := repo.ctx.DB.Model(QuestionEntity.Answer{})
	if by.QuestionId != 0 {
		db.Where("question_id", by.QuestionId)
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

func (repo AnswerRepository) GetUserCarbon(by repotypes.GetQuestionUserCarbon) float64 {
	carbon := float64(0)
	db := repo.ctx.DB.Model(QuestionEntity.Answer{})
	if by.QuestionId != 0 {
		db.Where("question_id", by.QuestionId)
	}
	if by.Uid != 0 {
		db.Where("user_id", by.Uid)
	}
	db.Select("sum(carbon) as carbon")
	if err := db.Find(&carbon).Error; err != nil {
		panic(err)
	}
	return carbon
}
