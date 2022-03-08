package service

import (
	"github.com/pkg/errors"
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

func (u UserService) GetUserById(id int64) (*entity.User, error) {
	return u.r.GetUserById(id)
}
func (u UserService) GetUserByOpenId(openId string) (*entity.User, error) {
	user := u.r.GetUserBy(repository.GetUserBy{
		OpenId: openId,
	})
	if user.ID == 0 {
		return nil, errors.New("未查询到用户信息")
	}
	return &user, nil
}
func (u UserService) GetUserByToken(token string) (*entity.User, error) {
	var authUser auth.User
	err := util.ParseToken(token, &authUser)
	if err != nil {
		return nil, err
	}
	return u.r.GetUserById(authUser.Id)
}
func (u UserService) CreateUserToken(id int64) (string, error) {
	user, err := u.r.GetUserById(id)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("用户不存在")
	}
	return util.CreateToken(auth.User{
		Id:     user.ID,
		Mobile: user.PhoneNumber,
	})
}
