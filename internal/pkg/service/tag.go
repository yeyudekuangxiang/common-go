package service

import (
	"github.com/mlogclub/simple"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
)

type (
	TagService interface {
		List(cnq *simple.SqlCnd) (list []entity.Tag)
		GetTagPageList(param repository.GetTagPageListBy) ([]entity.Tag, int64, error)
		GetOne(id int64) (entity.Tag, error)
	}

	defaultTagService struct {
		ctx      *mioContext.MioContext
		tagModel repository.TagModel
	}
)

func NewTagService(ctx *mioContext.MioContext) TagService {
	return defaultTagService{
		ctx:      ctx,
		tagModel: repository.NewTagModel(ctx),
	}
}

func (srv defaultTagService) List(cnq *simple.SqlCnd) (list []entity.Tag) {
	return srv.tagModel.List(cnq)
}

func (srv defaultTagService) GetTagPageList(param repository.GetTagPageListBy) ([]entity.Tag, int64, error) {
	list, total := srv.tagModel.GetTagPageList(param)
	return list, total, nil
}

func (srv defaultTagService) GetOne(id int64) (entity.Tag, error) {
	tag := srv.tagModel.GetById(id)
	if tag.Id == 0 {
		return entity.Tag{}, errno.ErrCommon.WithMessage("未找到该标签")
	}
	return tag, nil
}
