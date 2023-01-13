package community

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

type (
	TagModel interface {
		List(cnd *simple.SqlCnd) (list []entity.Tag)
		GetTagPageList(by repository.GetTagPageListBy) (list []entity.Tag, total int64)
		GetById(id int64) entity.Tag
		Delete(id int64) error
		Update(tag *entity.Tag) error
		Create(tag *entity.Tag) error
	}

	defaultTagModel struct {
		ctx *context.MioContext
	}
)

func NewTagModel(ctx *context.MioContext) TagModel {
	return defaultTagModel{
		ctx: ctx,
	}
}

func (m defaultTagModel) Delete(id int64) error {
	return m.ctx.DB.Delete(&entity.Tag{}, id).Error
}

func (m defaultTagModel) Update(tag *entity.Tag) error {
	return m.ctx.DB.Updates(tag).Error
}

func (m defaultTagModel) Create(tag *entity.Tag) error {
	return m.ctx.DB.Save(tag).Error
}

func (m defaultTagModel) GetById(id int64) entity.Tag {
	tag := entity.Tag{}
	err := m.ctx.DB.Model(entity.Tag{}).Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return tag
}

func (m defaultTagModel) List(cnd *simple.SqlCnd) (list []entity.Tag) {
	cnd.Find(app.DB, &list)
	return
}

func (m defaultTagModel) GetTagPageList(by repository.GetTagPageListBy) (list []entity.Tag, total int64) {
	list = make([]entity.Tag, 0)
	db := m.ctx.DB.Model(entity.Tag{})
	if by.ID != 0 {
		db.Where("id = ?", by.ID)
	}

	if by.ID != 0 {
		db.Where("id = ?", by.ID)
	}
	if by.Description != "" {
		db.Where("description like ?", by.Description+"%")
	}
	if by.Name != "" {
		db.Where("name = ?", by.Name)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case entity.OrderByTagSortAsc:
			db.Order("sort asc")
		case entity.OrderByTagSortDesc:
			db.Order("sort desc")
		}
	}

	err := db.Count(&total).Offset(by.Offset).Limit(by.Limit).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
