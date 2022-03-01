package repository

import (
	"github.com/mlogclub/simple"
	"mio/core/app"
	"mio/model"
)

var DefaultTopicRepository ITopicRepository = NewTopicRepository()

type ITopicRepository interface {
	List(cnd *simple.SqlCnd) (list []model.Topic)
}

func NewTopicRepository() TopicRepository {
	return TopicRepository{}
}

type TopicRepository struct {
}

func (u TopicRepository) List(cnd *simple.SqlCnd) (list []model.Topic) {
	cnd.Find(app.DB, &list)
	return
}
