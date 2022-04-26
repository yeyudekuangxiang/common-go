package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultOaAuthWhiteRepository = OaAuthWhiteRepository{DB: app.DB}

type OaAuthWhiteRepository struct {
	DB *gorm.DB
}

func (repo OaAuthWhiteRepository) FindBy(by FindOaAuthWhiteBy) entity.OaAuthWhite {
	white := entity.OaAuthWhite{}
	db := repo.DB.Model(white)
	if by.AppId != "" {
		db.Where("appid = ?", by.AppId)
	}
	if by.Domain != "" {
		db.Where("domain = ?", by.Domain)
	}

	err := db.First(&white).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return white
}
