package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultSystemAdminRepository ISystemAdminRepository = NewSystemAdminRepository(app.DB)

type ISystemAdminRepository interface {
	// GetAdminById 根据管理员id获取管理员信息
	GetAdminById(int) entity.SystemAdmin
	FindAdminBy(by FindAdminBy) entity.SystemAdmin
	GetAdminList(param GetAdminListBy) []entity.SystemAdmin
}

func NewSystemAdminRepository(db *gorm.DB) SystemAdminRepository {
	return SystemAdminRepository{DB: db}
}

type SystemAdminRepository struct {
	DB *gorm.DB
}

func (repo SystemAdminRepository) GetAdminById(id int) entity.SystemAdmin {
	var admin entity.SystemAdmin
	if err := repo.DB.First(&admin, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return admin
}
func (repo SystemAdminRepository) Create(admin *entity.SystemAdmin) error {
	return repo.DB.Create(admin).Error
}
func (repo SystemAdminRepository) Save(admin *entity.SystemAdmin) error {
	return repo.DB.Save(admin).Error
}
func (repo SystemAdminRepository) GetAdminList(param GetAdminListBy) []entity.SystemAdmin {
	list := make([]entity.SystemAdmin, 0)
	err := repo.DB.Find(&list).Error
	if err != nil {
		panic(err)
	}

	return list
}
func (repo SystemAdminRepository) FindAdminBy(by FindAdminBy) entity.SystemAdmin {
	admin := entity.SystemAdmin{}
	db := app.DB.Model(admin)
	if by.Account != "" {
		db.Where("account = ?", by.Account)
	}
	err := db.First(&admin).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return admin
}
