package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultBdSceneUserRepository = BdSceneUserRepository{DB: app.DB}

type BdSceneUserRepository struct {
	DB *gorm.DB
}

func (repo BdSceneUserRepository) FindByCh(ch string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := repo.DB.Where("ch = ?", ch).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}
func (repo BdSceneUserRepository) Create(data *entity.BdSceneUser) error {
	return repo.DB.Model(&entity.BdSceneUser{}).Save(data).Error
}
