package repository

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultTagRepository ITagRepository = NewTagRepository()

type ITagRepository interface {
	List(cnd *simple.SqlCnd) (list []entity.Tag)
	GetTagPageList(by GetTagPageListBy) (list []entity.Tag, total int64)
	GetById(id int64) entity.Tag
}

func NewTagRepository() TagRepository {
	return TagRepository{}
}

type TagRepository struct {
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
	if by.ID > 0 {
		db.Where("id = ?", by.ID)
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
