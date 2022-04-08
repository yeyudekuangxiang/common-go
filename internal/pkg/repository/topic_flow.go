package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultTopicFlowRepository = NewTopicFlowRepository(app.DB)

func NewTopicFlowRepository(db *gorm.DB) TopicFlowRepository {
	return TopicFlowRepository{DB: db}
}

type TopicFlowRepository struct {
	DB *gorm.DB
}

func (repo TopicFlowRepository) FindBy(by FindTopicFlowBy) entity.TopicFlow {

	flow := entity.TopicFlow{}

	db := repo.DB.Model(entity.TopicFlow{})
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
func (repo TopicFlowRepository) Save(flow *entity.TopicFlow) error {
	return repo.DB.Save(flow).Error
}
