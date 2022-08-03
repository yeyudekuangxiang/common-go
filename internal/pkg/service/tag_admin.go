package service

import (
	"github.com/mlogclub/simple"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultTagAdminService = NewTagAdminService(repository.DefaultTagRepository)

func NewTagAdminService(r repository.ITagRepository) TagAdminService {
	return TagAdminService{
		r: r,
	}
}

type TagAdminService struct {
	r repository.ITagRepository
}

func (u TagAdminService) List(cnq *simple.SqlCnd) (list []entity.Tag) {
	return u.r.List(cnq)
}

func (u TagAdminService) GetTagPageList(param repository.GetTagPageListBy) ([]entity.Tag, int64, error) {
	list, total := u.r.GetTagPageList(param)
	return list, total, nil
}

func (u TagAdminService) Delete(id int64) error {
	if err := u.r.Delete(id); err != nil {
		return err
	}
	return nil
}

func (u TagAdminService) Update(tag repository.UpdateTag) error {
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
	if err := u.r.Update(tagModel); err != nil {
		return err
	}
	return nil
}

func (u TagAdminService) Detail(id int64) (*entity.Tag, error) {
	detail, err := u.r.Detail(id)
	if err != nil {
		return &entity.Tag{}, err
	}
	return detail, nil
}

func (u TagAdminService) Create(tag repository.CreateTag) error {
	tagModel := &entity.Tag{
		Name:        tag.Name,
		Description: tag.Description,
		Img:         tag.Image,
		//Icon:        "",
	}
	if err := u.r.Create(tagModel); err != nil {
		return err
	}
	return nil
}
