package repository

import (
	"gorm.io/gorm"
	"mio/core/app"
	"mio/model/entity"
)

var DefaultTopicUserFlowRepository = NewTopicUserFlowRepository(app.DB)

func NewTopicUserFlowRepository(db *gorm.DB) TopicUserFlowRepository {
	return TopicUserFlowRepository{DB: db}
}

type TopicUserFlowRepository struct {
	DB *gorm.DB
}

func (repo TopicUserFlowRepository) FindBy(by FindTopicFlowBy) entity.TopicUserFlow {

	flow := entity.TopicUserFlow{}

	db := repo.DB.Model(entity.TopicUserFlow{})
	if by.TopicId > 0 {
		db.Where("topic_id = ?", by.TopicId)
	}
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if err := db.First(&flow).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return flow
}
func (repo TopicUserFlowRepository) Save(flow *entity.TopicUserFlow) error {
	return repo.DB.Save(flow).Error
}
