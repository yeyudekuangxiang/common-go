package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	QuestionEntity "mio/internal/pkg/model/entity/question"
	"mio/internal/pkg/repository/repotypes"
)

func NewIndexIconRepository(ctx *context.MioContext) *IndexIconRepository {
	return &IndexIconRepository{
		ctx: ctx,
	}
}

type IndexIconRepository struct {
	ctx *context.MioContext
}

func (repo IndexIconRepository) Save(transaction *entity.IndexIcon) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo IndexIconRepository) CreateInBatches(transaction []entity.IndexIcon) error {
	return repo.ctx.DB.CreateInBatches(transaction, len(transaction)).Error
}

func (repo IndexIconRepository) Delete(by *repotypes.DeleteIndexIconDO) error {
	db := repo.ctx.DB.Model(entity.IndexIcon{})
	if by.Id == 0 {
		return errors.New("id不能为空")
	}
	db.Where("id", by.Id)
	return db.Updates(by).Error
}

func (repo IndexIconRepository) GetOne(do repotypes.GetBannerOneDO) (*entity.IndexIcon, bool, error) {
	data := entity.IndexIcon{}
	db := repo.ctx.DB.Model(data)
	if do.ID != 0 {
		db.Where("id = ?", do.ID)
	}
	err := db.First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &data, true, nil
}

func (repo IndexIconRepository) GetListBy(by repotypes.GetQuestOptionGetListBy) ([]entity.IndexIcon, error) {
	list := make([]entity.IndexIcon, 0)
	db := repo.ctx.DB.Model(entity.IndexIcon{})
	if len(by.SubjectIds) != 0 {
		db.Where("subject_id in (?) ", by.SubjectIds)
	}
	db.Order("sort desc")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}

func (repo IndexIconRepository) GetListByUid(by repotypes.GetQuestionOptionGetListByUid) ([]QuestionEntity.Answer, error) {
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
