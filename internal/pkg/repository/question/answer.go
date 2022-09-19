package question

import (
	"github.com/pkg/errors"
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

func (repo AnswerRepository) Delete(by *repotypes.DeleteQuestionAnswerDO) error {
	db := repo.ctx.DB.Model(QuestionEntity.Answer{})
	if by.Uid == 0 {
		return errors.New("用户id不能为空")
	}
	db.Where("user_id", by.Uid)
	if by.QuestionId != 0 {
		db.Where("question_id", by.QuestionId)
	}
	return db.Updates(by).Error
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
	sum := struct {
		Carbon float64
	}{}
	db := repo.ctx.DB.Model(QuestionEntity.Answer{})
	if by.QuestionId != 0 {
		db.Where("question_id", by.QuestionId)
	}
	if by.Uid != 0 {
		db.Where("user_id", by.Uid)
	}
	db.Select("sum(carbon) as carbon")
	if err := db.Take(&sum).Error; err != nil {
		panic(err)
	}
	return sum.Carbon
}

func (repo AnswerRepository) GetUserAnswer(by repotypes.GetQuestionUserCarbon) []repotypes.UserAnswerStruct {
	var list []repotypes.UserAnswerStruct
	db := repo.ctx.DB.Model(QuestionEntity.Answer{})

	db.Joins("left join question_subject on question_subject.subject_id = question_answer.subject_id")
	if by.QuestionId != 0 {
		db.Where("question_answer.question_id", by.QuestionId)
	}
	if by.Uid != 0 {
		db.Where("question_answer.user_id", by.Uid)
	}
	db.Where("is_delete", 0)
	db.Select("category_id,sum(carbon) as carbon")
	db.Group("category_id")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
