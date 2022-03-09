package service

import (
	"github.com/pkg/errors"
	"mio/internal/util"
	"mio/model"
	"mio/model/auth"
	"mio/model/entity"
	"mio/repository"
	"time"
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
	if id == 0 {
		return &entity.User{}, nil
	}
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
		Id:        user.ID,
		Mobile:    user.PhoneNumber,
		CreatedAt: model.Time{Time: time.Now()},
	})
}
func (u UserService) CreateUser(param CreateUserParam) (*entity.User, error) {
	user := entity.User{}
	if err := util.MapTo(param, &user); err != nil {
		return nil, err
	}
	user.Time = model.NewTime()
	return &user, repository.DefaultUserRepository.Save(&user)
}
func (u UserService) GetUserBy(by repository.GetUserBy) entity.User {
	return repository.DefaultUserRepository.GetUserBy(by)
}
func (u UserService) FindOrCreateByMobile(mobile string) (*entity.User, error) {
	user := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{
		Mobile: mobile,
		Source: entity.UserSourceMobile,
	})

	if user.ID > 0 {
		return &user, nil
	}
	return u.CreateUser(CreateUserParam{
		OpenId:      mobile,
		Nickname:    "手机用户" + mobile[len(mobile)-4:],
		PhoneNumber: mobile,
		Source:      entity.UserSourceMobile,
		UnionId:     mobile,
	})
}

// FindUserBySource 根据用户id 获取指定平台的用户
func (u UserService) FindUserBySource(source entity.UserSource, userId int64) (*entity.User, error) {
	user, err := repository.DefaultUserRepository.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil || user.PhoneNumber == "" {
		return &entity.User{}, nil
	}

	sourceUer := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{
		Mobile: user.PhoneNumber,
		Source: source,
	})

	return &sourceUer, nil
}
