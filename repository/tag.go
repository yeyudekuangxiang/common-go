package repository

import (
	"github.com/mlogclub/simple"
	"mio/core/app"
	"mio/model"
)

var DefaultTagRepository ITagRepository = NewTagRepository()

type ITagRepository interface {
	List(cnd *simple.SqlCnd) (list []model.Tag)
}

func NewTagRepository() TagRepository {
	return TagRepository{}
}

type TagRepository struct {
}

func (u TagRepository) List(cnd *simple.SqlCnd) (list []model.Tag) {
	cnd.Find(app.DB, &list)
	return
}
