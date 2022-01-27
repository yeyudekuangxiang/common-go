package repository

import (
	"mio/core/app"
	"mio/model/entity"
	"gorm.io/gorm"
)

var DefaultUserRepository IUserRepository = NewUserRepository()

type IUserRepository interface {
	//根据guid获取用户信息
	GetUserByGuid(string) (*entity.User, error)
	//根据用id获取用户信息
	GetUserById(int) (*entity.User, error)
}

func NewUserRepository() UserRepository {
	return UserRepository{}
}

type UserRepository struct {
}

func (u UserRepository) GetUserById(id int) (*entity.User, error) {
	var user entity.User
	if err := app.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u UserRepository) GetUserByGuid(guid string) (*entity.User, error) {
	var user entity.User
	if err := app.DB.Where("guid = ?", guid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
