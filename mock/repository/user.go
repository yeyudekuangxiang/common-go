package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

func NewUserMockRepository() UserMockRepository {
	return UserMockRepository{}
}

type UserMockRepository struct {
	db *gorm.DB
}

func (u UserMockRepository) Save(user *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserMockRepository) GetUserById(i int64) entity.User {
	//TODO implement me
	panic("implement me")
}

func (u UserMockRepository) GetShortUserListBy(by repository.GetUserListBy) []entity.ShortUser {
	//TODO implement me
	panic("implement me")
}

func (u UserMockRepository) GetUserListBy(by repository.GetUserListBy) []entity.User {
	//TODO implement me
	panic("implement me")
}

func (u UserMockRepository) GetUserBy(by repository.GetUserBy) entity.User {
	//TODO implement me
	panic("implement me")
}

func (u UserMockRepository) GetShortUserBy(by repository.GetUserBy) entity.ShortUser {
	//TODO implement me
	panic("implement me")
}

func (u UserMockRepository) GetUserByGuid(guid string) (*entity.User, error) {
	return &entity.User{
		ID:       1,
		Nickname: fmt.Sprintf("mock%s", guid),
	}, nil
}
