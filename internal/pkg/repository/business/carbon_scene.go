package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCarbonSceneRepository = CarbonSceneRepository{DB: app.DB}

type CarbonSceneRepository struct {
	DB *gorm.DB
}

func (repo CarbonSceneRepository) FindScene(t business.CarbonType) business.CarbonScene {
	scene := business.CarbonScene{}
	err := repo.DB.Model(scene).Where("type = ?", t).Take(&scene).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return scene
}

type ICarbonSceneRepository interface {
	Save(CarbonScene *business.CarbonScene) error
}

func (repo CarbonSceneRepository) GetCarbonSceneListBy(by GetCarbonSceneListBy) []business.CarbonScene {
	list := make([]business.CarbonScene, 0)
	db := app.DB.Model(business.CarbonScene{})
	if len(by.Ids) > 0 {
		db.Where("id in (?)", by.Ids)
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}

	return list
}
