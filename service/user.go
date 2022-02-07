package service

import (
	"mio/internal/util"
	"mio/model/auth"
	"mio/model/entity"
	"mio/repository"
)

var DefaultUserService = NewUserService(repository.DefaultUserRepository)

func NewUserService(r repository.IUserRepository) UserService {
	return UserService{
		r: r,
	}
}

type UserService struct {
	r repository.IUserRepository
}

func (u UserService) GetUserById(id int) (*entity.User, error) {
	return u.r.GetUserById(id)
}

func (u UserService) GetUserByToken(token string) (*entity.User, error) {
	var authUser auth.User
	err := util.ParseToken(token, &authUser)
	if err != nil {
		return nil, err
	}
	return u.r.GetUserByGuid(authUser.Guid)
}
