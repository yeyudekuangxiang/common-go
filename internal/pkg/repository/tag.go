package repository

import (
	"github.com/mlogclub/simple"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultTagRepository ITagRepository = NewTagRepository()

type ITagRepository interface {
	List(cnd *simple.SqlCnd) (list []entity.Tag)
	GetTagPageList(by GetTagPageListBy) (list []entity.Tag, total int64)
}

func NewTagRepository() TagRepository {
	return TagRepository{}
}

type TagRepository struct {
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
