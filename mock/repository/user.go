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

func (u UserMockRepository) GetUserById(i int64) (*entity.User, error) {
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
