package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultBdSceneRepository = BdSceneRepository{DB: app.DB}

type BdSceneRepository struct {
	DB *gorm.DB
}

func (repo BdSceneRepository) FindByCh(ch string) entity.BdScene {
	item := entity.BdScene{}
	err := repo.DB.Where("ch = ?", ch).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (repo BdSceneRepository) SceneToType(ch string) entity.PointTransactionType {
	switch ch {
	case "lvmiao":
		return entity.POINT_ECAR
	case "jinhuaxing":
		return entity.POINT_JHX
	}
	return entity.POINT_ECAR
}