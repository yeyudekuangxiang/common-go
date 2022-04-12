package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultAdminRepository IAdminRepository = NewAdminRepository()

type IAdminRepository interface {
	// GetAdminById 根据管理员id获取管理员信息
	GetAdminById(int) entity.Admin
}

func NewAdminRepository() AdminRepository {
	return AdminRepository{}
}

type AdminRepository struct {
}

func (a AdminRepository) GetAdminById(id int) entity.Admin {
	var admin entity.Admin
	if err := app.DB.First(&admin, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}

	}
	return admin
}
