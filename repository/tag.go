package repository

import (
	"github.com/mlogclub/simple"
	"mio/core/app"
	"mio/model/entity"
)

var DefaultTagRepository ITagRepository = NewTagRepository()

type ITagRepository interface {
	List(cnd *simple.SqlCnd) (list []entity.Tag)
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
