package kumiaoCommunity

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	CommunityActivitiesTagModel interface {
		List(params GetActivitiesTagListParams) ([]entity.CommunityActivitiesTag, error)
		GetPageList(params GetActivitiesTagPageListParams) ([]entity.CommunityActivitiesTag, int64, error)
		GetById(id int64) (entity.CommunityActivitiesTag, error)
		Delete(id int64) error
		Update(tag *entity.CommunityActivitiesTag) error
		Create(tag *entity.CommunityActivitiesTag) error
	}

	defaultCommunityActivitiesTagModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultCommunityActivitiesTagModel) List(params GetActivitiesTagListParams) ([]entity.CommunityActivitiesTag, error) {
	var list []entity.CommunityActivitiesTag
	db := d.ctx.DB.Model(&entity.CommunityActivitiesTag{})
	err := db.Find(&list).Error
	if err != nil {
		return []entity.CommunityActivitiesTag{}, err
	}
	return list, nil
}

func (d defaultCommunityActivitiesTagModel) GetPageList(params GetActivitiesTagPageListParams) ([]entity.CommunityActivitiesTag, int64, error) {
	list := make([]entity.CommunityActivitiesTag, 0)
	var total int64
	db := d.ctx.DB.Model(entity.CommunityActivitiesTag{})
	if params.ID != 0 {
		db.Where("id = ?", params.ID)
	}
	if params.ID != 0 {
		db.Where("id = ?", params.ID)
	}
	if params.Description != "" {
		db.Where("description like ?", params.Description+"%")
	}
	if params.Name != "" {
		db.Where("name = ?", params.Name)
	}

	err := db.Count(&total).Offset(params.Offset).Limit(params.Limit).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (d defaultCommunityActivitiesTagModel) GetById(id int64) (entity.CommunityActivitiesTag, error) {
	var resp entity.CommunityActivitiesTag
	err := d.ctx.DB.Model(&resp).WithContext(d.ctx.Context).First(&resp, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.CommunityActivitiesTag{}, nil
		}
		return entity.CommunityActivitiesTag{}, err
	}
	return resp, nil
}

func (d defaultCommunityActivitiesTagModel) Delete(id int64) error {
	if err := d.ctx.DB.WithContext(d.ctx.Context).Delete(&entity.CommunityActivitiesTag{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (d defaultCommunityActivitiesTagModel) Update(tag *entity.CommunityActivitiesTag) error {
	if err := d.ctx.DB.Model(&entity.CommunityActivitiesTag{}).WithContext(d.ctx.Context).Save(tag).Error; err != nil {
		return err
	}
	return nil
}

func (d defaultCommunityActivitiesTagModel) Create(tag *entity.CommunityActivitiesTag) error {
	if err := d.ctx.DB.Model(&entity.CommunityActivitiesTag{}).WithContext(d.ctx.Context).Save(tag).Error; err != nil {
		return err
	}
	return nil
}

func NewCommunityActivitiesTagModel(ctx *mioContext.MioContext) CommunityActivitiesTagModel {
	return defaultCommunityActivitiesTagModel{
		ctx: ctx,
	}
}
