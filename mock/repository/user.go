package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mio/model/entity"
	"mio/repository"
)

func NewUserMockRepository() UserMockRepository {
	return UserMockRepository{}
}

type UserMockRepository struct {
	db *gorm.DB
}

func (u UserMockRepository) GetUserBy(by repository.GetUserBy) entity.User {
	//TODO implement me
	panic("implement me")
}

func (u UserMockRepository) GetShortUserBy(by repository.GetUserBy) entity.ShortUser {
	//TODO implement me
	panic("implement me")
}

func (u UserMockRepository) GetUserById(id int) (*entity.User, error) {
	return &entity.User{
		ID:       int64(id),
		Nickname: fmt.Sprintf("mock%d", id),
	}, nil
}

func (u UserMockRepository) GetUserByGuid(guid string) (*entity.User, error) {
	return &entity.User{
		ID:       1,
		Nickname: fmt.Sprintf("mock%s", guid),
	}, nil
}
