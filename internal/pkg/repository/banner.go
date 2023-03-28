package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

var DefaultBannerRepository = BannerRepository{DB: app.DB}

type BannerRepository struct {
	DB *gorm.DB
}

func (repo BannerRepository) GetById(id int64) entity.Banner {
	banner := entity.Banner{}
	err := repo.DB.Model(banner).First(&banner, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return banner
}

func (repo BannerRepository) GetExistOne(do repotypes.GetBannerExistDO) (entity.Banner, error) {
	banner := entity.Banner{}
	db := repo.DB.Model(banner)
	if do.Name != "" {
		db.Where("name = ?", do.Name)
	}
	if do.ImageUrl != "" {
		db.Or("image_url = ?", do.ImageUrl)
	}
	if do.NotId != 0 {
		db.Not("id", do.NotId)
	}
	err := db.First(&banner).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return banner, err
	}
	return banner, nil
}

func (repo BannerRepository) GetOne(do repotypes.GetBannerOneDO) (entity.Banner, error) {
	banner := entity.Banner{}
	db := repo.DB.Model(banner)
	if do.ID != 0 {
		db.Where("id = ?", do.ID)
	}
	err := db.First(&banner).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return banner, err
	}
	return banner, nil
}
func (repo BannerRepository) List(do repotypes.GetBannerListDO) ([]entity.Banner, error) {
	db := repo.DB.Model(entity.Banner{})
	if do.Scene != "" {
		db.Where("scene = ?", do.Scene)
	}
	if do.Type != "" {
		db.Where("type = ?", do.Type)
	}
	if do.Status != 0 {
		db.Where("status = ?", do.Status)
	}

	for _, orderBy := range do.OrderBy {
		switch orderBy {
		case entity.OrderByBannerSortAsc:
			db.Order("sort asc")
		}
	}
	list := make([]entity.Banner, 0)
	return list, db.Find(&list).Error
}
func (repo BannerRepository) Page(do repotypes.GetBannerPageDO) ([]entity.Banner, int64, error) {
	db := repo.DB.Model(entity.Banner{})
	if do.Scene != "" {
		db.Where("scene = ?", do.Scene)
	}
	if do.Type != "" {
		db.Where("type = ?", do.Type)
	}
	if do.Status != 0 {
		db.Where("status = ?", do.Status)
	}
	if do.Name != "" {
		db.Where("name like ?", "%"+do.Name+"%")
	}
	if len(do.Displays) != 0 {
		db.Where("display in (?)", do.Displays)
	}
	for _, orderBy := range do.OrderBy {
		switch orderBy {
		case entity.OrderByBannerSortAsc:
			db.Order("sort asc")
		}
	}
	list := make([]entity.Banner, 0)
	var total int64
	return list, total, db.Count(&total).Offset(do.Offset).Limit(do.Limit).Find(&list).Error
}
func (repo BannerRepository) Save(banner *entity.Banner) error {
	return repo.DB.Save(banner).Error
}
func (repo BannerRepository) Create(do *entity.Banner) error {
	return repo.DB.Create(do).Error
}
