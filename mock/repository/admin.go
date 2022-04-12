package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/model/entity"
)

func NewAdminMockRepository() AdminMockRepository {
	return AdminMockRepository{}
}

type AdminMockRepository struct {
	db *gorm.DB
}

func (a AdminMockRepository) GetAdminById(id int) entity.Admin {
	return entity.Admin{
		ID:       id,
		UName:    fmt.Sprintf("mock%d", id),
		RealName: fmt.Sprintf("mock%d", id),
	}
}
