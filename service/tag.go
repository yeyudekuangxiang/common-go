package service

import (
	"github.com/mlogclub/simple"
	"mio/model/entity"
	"mio/repository"
)

var DefaultTagService = NewTagService(repository.DefaultTagRepository)

func NewTagService(r repository.ITagRepository) TagService {
	return TagService{
		r: r,
	}
}

type TagService struct {
	r repository.ITagRepository
}

func (u TagService) List(cnq *simple.SqlCnd) (list []entity.Tag) {
	return u.r.List(cnq)
}
