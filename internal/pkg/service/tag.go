package service

import (
	"github.com/mlogclub/simple"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
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

func (u TagService) GetOne(id int64) (entity.Tag, error) {
	tag := u.r.GetById(id)
	if tag.Id == 0 {
		return entity.Tag{}, errno.ErrCommon.WithMessage("未找到该标签")
	}
	return tag, nil
}
