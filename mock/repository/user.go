package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mio/model/entity"
)

func NewUserMockRepository() UserMockRepository {
	return UserMockRepository{}
}

type UserMockRepository struct {
	db *gorm.DB
}

func (u UserMockRepository) GetUserById(id int) (*entity.User, error) {
	return &entity.User{
		ID:       id,
		Guid:     "13123",
		Nickname: fmt.Sprintf("mock%d", id),
	}, nil
}

func (u UserMockRepository) GetUserByGuid(guid string) (*entity.User, error) {
	return &entity.User{
		ID:       1,
		Guid:     guid,
		Nickname: fmt.Sprintf("mock%s", guid),
	}, nil
}
