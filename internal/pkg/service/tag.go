package service

import (
	"github.com/mlogclub/simple"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
)

var DefaultTagService = NewTagService(repository2.DefaultTagRepository)

func NewTagService(r repository2.ITagRepository) TagService {
	return TagService{
		r: r,
	}
}

type TagService struct {
	r repository2.ITagRepository
}

func (u TagService) List(cnq *simple.SqlCnd) (list []entity.Tag) {
	return u.r.List(cnq)
}

func (u TagService) GetTagPageList(param repository2.GetTagPageListBy) ([]entity.Tag, int64, error) {
	list, total := u.r.GetTagPageList(param)
	return list, total, nil
}
