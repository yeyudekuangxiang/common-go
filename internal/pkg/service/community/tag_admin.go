package community

import (
	"github.com/mlogclub/simple"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

type (
	TagAdminService interface {
		List(cnq *simple.SqlCnd) (list []entity.Tag)
		GetTagPageList(param repository.GetTagPageListBy) ([]entity.Tag, int64, error)
		Delete(id int64) error
		Update(tag repository.UpdateTag) error
		Detail(id int64) entity.Tag
		Create(tag repository.CreateTag) error
	}

	defaultTagAdminService struct {
		ctx      *mioContext.MioContext
		tagModel repository.TagModel
	}
)

func NewTagAdminService(ctx *mioContext.MioContext) TagAdminService {
	return defaultTagAdminService{
		ctx:      ctx,
		tagModel: repository.NewTagModel(ctx),
	}
}

func (srv defaultTagAdminService) List(cnq *simple.SqlCnd) (list []entity.Tag) {
	return srv.tagModel.List(cnq)
}

func (srv defaultTagAdminService) GetTagPageList(param repository.GetTagPageListBy) ([]entity.Tag, int64, error) {
	list, total := srv.tagModel.GetTagPageList(param)
	return list, total, nil
}

func (srv defaultTagAdminService) Delete(id int64) error {
	if err := srv.tagModel.Delete(id); err != nil {
		return err
	}
	return nil
}

func (srv defaultTagAdminService) Update(tag repository.UpdateTag) error {
	tagModel := &entity.Tag{Id: tag.ID}
	if tag.Image != "" {
		tagModel.Img = tag.Image
	}
	if tag.Name != "" {
		tagModel.Name = tag.Name
	}
	if tag.Description != "" {
		tagModel.Description = tag.Description
	}
	if err := srv.tagModel.Update(tagModel); err != nil {
		return err
	}
	return nil
}

func (srv defaultTagAdminService) Detail(id int64) entity.Tag {
	return srv.tagModel.GetById(id)
}

func (srv defaultTagAdminService) Create(tag repository.CreateTag) error {
	tagModel := &entity.Tag{
		Name:        tag.Name,
		Description: tag.Description,
		Img:         tag.Image,
		//Icon:        "",
	}
	if err := srv.tagModel.Create(tagModel); err != nil {
		return err
	}
	return nil
}
