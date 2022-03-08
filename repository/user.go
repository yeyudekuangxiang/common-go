package repository

import (
	"gorm.io/gorm"
	"mio/core/app"
	"mio/model/entity"
)

var DefaultUserRepository IUserRepository = NewUserRepository()

type IUserRepository interface {
	// GetUserById 根据用id获取用户信息
	GetUserById(int64) (*entity.User, error)
	GetUserBy(by GetUserBy) entity.User
	GetShortUserBy(by GetUserBy) entity.ShortUser
	GetUserListBy(by GetUserListBy) []entity.User
}

func NewUserRepository() UserRepository {
	return UserRepository{}
}

type UserRepository struct {
}

func (u UserRepository) GetUserById(id int64) (*entity.User, error) {
	var user entity.User
	if err := app.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (u UserRepository) GetUserBy(by GetUserBy) entity.User {
	user := entity.User{}
	db := app.DB.Model(user)
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if err := db.First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return user
}
func (u UserRepository) GetShortUserBy(by GetUserBy) entity.ShortUser {
	user := entity.ShortUser{}
	db := app.DB.Model(entity.User{})
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if err := db.First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return user
}
func (u UserRepository) GetUserListBy(by GetUserListBy) []entity.User {
	list := make([]entity.User, 0)
	db := app.DB.Model(entity.User{})

	if by.Mobile != "" {
		db.Where("phone_number = ?", by.Mobile)
	}
	if len(by.UserIds) > 0 {
		db.Where("id in (?)", by.UserIds)
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
