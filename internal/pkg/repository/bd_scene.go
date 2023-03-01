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
	case "fastElectricity":
		return entity.POINT_FAST_ELECTRICITY
	case "yitongxing":
		return entity.POINT_YTX
	case "ykc":
		return entity.POINT_YKC
	}
	return entity.POINT_ECAR
}

func (repo BdSceneRepository) SceneToCarbonType(ch string) entity.CarbonTransactionType {
	switch ch {
	case "lvmiao":
		return entity.CARBON_ECAR
	case "jinhuaxing":
		return entity.CARBON_JHX
	case "yitongxing":
		return entity.CARBON_YTX
	case "fastElectricity":
		return entity.CARBON_FAST_ELECTRICITY
	case "ykc":
		return entity.CARBON_YKC
	}

	return ""
}
