package business

import (
	"errors"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/model/entity/business"
	brepo "mio/internal/pkg/repository/business"
	util2 "mio/internal/pkg/util"
	"time"
)

var DefaultUserService = UserService{repo: brepo.DefaultUserRepository}

type UserService struct {
	repo brepo.UserRepository
}

func (u UserService) GetBusinessUserBy(param GetBusinessUserParam) (*business.User, error) {
	user := u.repo.GetUserBy(brepo.GetUserBy{Mobile: param.Mobile})
	if user.ID != 0 {
		return &user, nil
	}
	return nil, errors.New("非企业用户,请联系管理员开通企业版权限")
}

// GetBusinessUserById 根据用户id查询用户信息
func (u UserService) GetBusinessUserById(id int64) (*business.User, error) {
	user := u.repo.GetUserBy(brepo.GetUserBy{ID: id})
	if user.ID != 0 {
		return &user, nil
	}
	return nil, errors.New("非企业用户,请联系管理员开通企业版权限")
}

// GetBusinessUserByUid 根据用户uid查询用户信息
func (u UserService) GetBusinessUserByUid(uid string) (*business.User, error) {
	user := u.repo.GetUserBy(brepo.GetUserBy{Uid: uid})
	if user.ID != 0 {
		return &user, nil
	}
	return nil, errors.New("非企业用户,请联系管理员开通企业版权限")
}

// GetBusinessUserByIds 批量查询用户信息
func (u UserService) GetBusinessUserByIds(ids []int64) []business.User {
	users := u.repo.GetUserListBy(brepo.GetUserListBy{Ids: ids})
	return users
}

func (u UserService) CreateBusinessUserToken(user *business.User) (string, error) {
	return util2.CreateToken(auth.BusinessUser{
		ID:        user.ID,
		Mobile:    user.Mobile,
		CreatedAt: model.Time{Time: time.Now()},
	})
}

func (u UserService) GetBusinessUserByToken(token string) (*business.User, error) {
	var user auth.BusinessUser
	err := util2.ParseToken(token, &user)
	if err != nil {
		return nil, err
	}
	newInfo := u.repo.GetUserBy(brepo.GetUserBy{ID: user.ID})
	if newInfo.ID > 0 {
		return &newInfo, nil
	}
	return nil, nil
}
