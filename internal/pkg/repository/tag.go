package repository

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultTagRepository ITagRepository = NewTagRepository(app.DB)

type ITagRepository interface {
	List(cnd *simple.SqlCnd) (list []entity.Tag)
	GetTagPageList(by GetTagPageListBy) (list []entity.Tag, total int64)
	GetById(id int64) entity.Tag
	Delete(id int64) error
	Update(tag *entity.Tag) error
	Create(tag *entity.Tag) error
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return TagRepository{DB: db}
}

type TagRepository struct {
	DB *gorm.DB
}

func (u TagRepository) Delete(id int64) error {
	return u.DB.Delete(&entity.Tag{}, id).Error
}

func (u TagRepository) Update(tag *entity.Tag) error {
	return u.DB.Updates(tag).Error
}

func (u TagRepository) Create(tag *entity.Tag) error {
	return u.DB.Save(tag).Error
}

func (u TagRepository) GetById(id int64) entity.Tag {
	tag := entity.Tag{}
	err := app.DB.Model(entity.Tag{}).Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return tag
}

func (u TagRepository) List(cnd *simple.SqlCnd) (list []entity.Tag) {
	cnd.Find(app.DB, &list)
	return
}

func (u TagRepository) GetTagPageList(by GetTagPageListBy) (list []entity.Tag, total int64) {
	list = make([]entity.Tag, 0)
	db := app.DB.Model(entity.Tag{})
	if by.Description != "" {
		db.Where("description like ?", by.Description+"%")
	}
	if by.Name != "" {
		db.Where("name = ?", by.Name)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case entity.OrderByTagSortAsc:
			db.Order("tag asc")
		case entity.OrderByTagSortDesc:
			db.Order("tag desc")
		}
	}

	err := db.Count(&total).Offset(by.Offset).Limit(by.Limit).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
