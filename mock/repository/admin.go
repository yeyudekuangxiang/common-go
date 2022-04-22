package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

func NewAdminMockRepository() AdminMockRepository {
	return AdminMockRepository{}
}

type AdminMockRepository struct {
	db *gorm.DB
}

func (a AdminMockRepository) FindAdminBy(by repository.FindAdminBy) entity.SystemAdmin {
	//TODO implement me
	panic("implement me")
}

func (a AdminMockRepository) GetAdminList(param repository.GetAdminListBy) []entity.SystemAdmin {
	//TODO implement me
	panic("implement me")
}

func (a AdminMockRepository) GetAdminById(id int) entity.SystemAdmin {
	return entity.SystemAdmin{
		ID:       id,
		Nickname: fmt.Sprintf("mock%d", id),
		RealName: fmt.Sprintf("mock%d", id),
	}
}
