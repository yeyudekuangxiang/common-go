package service

import (
	"github.com/mlogclub/simple"
	"mio/model"
	"mio/repository"
)

var DefaultTopicService = NewTopicService(repository.DefaultTopicRepository)

func NewTopicService(r repository.ITopicRepository) TopicService {
	return TopicService{
		r: r,
	}
}

type TopicService struct {
	r repository.ITopicRepository
}

func (u TopicService) List(cnq *simple.SqlCnd) (list []model.Topic) {
	return u.r.List(cnq)
}
