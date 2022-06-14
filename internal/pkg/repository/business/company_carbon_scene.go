package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCompanyCarbonSceneRepository = CompanyCarbonSceneRepository{DB: app.DB}

type CompanyCarbonSceneRepository struct {
	DB *gorm.DB
}
type ICompanyCarbonSceneRepository interface {
	Save(CompanyCarbonScene *business.CompanyCarbonScene) error
}

func (repo CompanyCarbonSceneRepository) FindCompanyScene(by FindCompanyCarbonSceneBy) business.CompanyCarbonScene {
	scene := business.CompanyCarbonScene{}
	db := repo.DB.Model(scene)
	if by.CompanyId != 0 {
		db.Where("b_company_id = ?", by.CompanyId)
	}
	if by.CarbonSceneId != 0 {
		db.Where("carbon_scene_id = ?", by.CarbonSceneId)
	}
	err := db.Take(&scene).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return scene
}

func (repo CompanyCarbonSceneRepository) GetCompanyCarbonSceneListBy(by GetCompanyCarbonSceneListBy) []business.CompanyCarbonScene {
	list := make([]business.CompanyCarbonScene, 0)
	db := app.DB.Model(business.CompanyCarbonScene{})
	if by.Status > 0 {
		db.Where("status = ?", by.Status)
	}

	if err := db.Find(&list).Order("sort desc").Error; err != nil {
		panic(err)
	}

	return list
}
