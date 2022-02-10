package repository

import (
	"fmt"
	"mio/model/entity"
	"gorm.io/gorm"
)

func NewAdminMockRepository() AdminMockRepository {
	return AdminMockRepository{}
}

type AdminMockRepository struct {
	db *gorm.DB
}

func (a AdminMockRepository) GetAdminById(id int) (*entity.Admin, error) {
	return &entity.Admin{
		ID:       id,
		UName:    fmt.Sprintf("mock%d", id),
		RealName: fmt.Sprintf("mock%d", id),
	}, nil
}
