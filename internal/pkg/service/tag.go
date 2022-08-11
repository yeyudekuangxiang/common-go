package service

import (
	"github.com/mlogclub/simple"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
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

func (u TagService) GetTagPageList(param repository.GetTagPageListBy) ([]entity.Tag, int64, error) {
	list, total := u.r.GetTagPageList(param)
	return list, total, nil
}
