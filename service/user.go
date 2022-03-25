package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"mio/config"
	"mio/core/app"
	"mio/internal/message"
	"mio/internal/util"
	"mio/model"
	"mio/model/auth"
	"mio/model/entity"
	"mio/repository"
	"strconv"
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
	user := u.r.GetUserById(id)
	return &user, nil
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
	user := u.r.GetUserById(authUser.Id)
	return &user, nil
}
func (u UserService) CreateUserToken(id int64) (string, error) {
	user := u.r.GetUserById(id)
	if user.ID == 0 {
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
	user := repository.DefaultUserRepository.GetUserById(userId)

	if user.ID == 0 || user.PhoneNumber == "" {
		return &entity.User{}, nil
	}

	sourceUer := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{
		Mobile: user.PhoneNumber,
		Source: source,
	})

	return &sourceUer, nil
}
func (u UserService) GetYZM(mobile string) (string, error) {
	code := ""
	for i := 0; i < 4; i++ {
		code += strconv.Itoa(rand.Intn(9))
	}
	//加入缓存
	cmd := app.Redis.Set(context.Background(), config.RedisKey.YZM+mobile, code, time.Second*30*60)
	fmt.Println(cmd.String())
	//发送短信
	message.SendYZM(mobile, code)

	return code, nil
}

func (u UserService) CheckYZM(mobile string, code string) bool {
	//取出缓存
	codeCmd := app.Redis.Get(context.Background(), config.RedisKey.YZM+mobile)
	fmt.Println(codeCmd.String())
	if codeCmd.Val() == code {
		fmt.Println("验证码验证通过")
		return true
	}

	return false
}
