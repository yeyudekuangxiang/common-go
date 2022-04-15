package service

import (
	"context"
	"fmt"
	"github.com/medivhzhan/weapp/v3/phonenumber"
	"github.com/pkg/errors"
	"math/rand"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
	util2 "mio/internal/pkg/util"
	"mio/internal/pkg/util/message"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultUserService = NewUserService(repository2.DefaultUserRepository)

func NewUserService(r repository2.IUserRepository) UserService {
	return UserService{
		r: r,
	}
}

type UserService struct {
	r repository2.IUserRepository
}

func (u UserService) GetUserById(id int64) (*entity.User, error) {
	if id == 0 {
		return &entity.User{}, nil
	}
	user := u.r.GetUserById(id)
	return &user, nil
}
func (u UserService) GetUserByOpenId(openId string) (*entity.User, error) {
	user := u.r.GetUserBy(repository2.GetUserBy{
		OpenId: openId,
	})
	if user.ID == 0 {
		return nil, errors.New("未查询到用户信息")
	}
	return &user, nil
}
func (u UserService) GetUserByToken(token string) (*entity.User, error) {
	var authUser auth.User
	err := util2.ParseToken(token, &authUser)
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
	return util2.CreateToken(auth.User{
		Id:        user.ID,
		Mobile:    user.PhoneNumber,
		CreatedAt: model.Time{Time: time.Now()},
	})
}

func (u UserService) CreateUser(param CreateUserParam) (*entity.User, error) {
	user := u.r.GetUserBy(repository2.GetUserBy{
		OpenId: param.OpenId,
		Source: param.Source,
	})
	if user.ID != 0 {
		return &user, nil
	}

	guid := ""
	if param.UnionId != "" {
		guid = u.r.GetGuid(param.UnionId)
		if guid == "" {
			guid = util2.UUID()
		}
	}

	user = entity.User{}
	if err := util2.MapTo(param, &user); err != nil {
		return nil, err
	}
	user.GUID = guid
	user.Time = model.NewTime()

	if param.UnionId != "" {
		app.DB.Where("unionid = ? and guid =''", param.UnionId).Update("guid", guid)
	}

	return &user, repository2.DefaultUserRepository.Save(&user)
}
func (u UserService) UpdateUserUnionId(id int64, unionid string) {
	if unionid == "" {
		return
	}

	user := u.r.GetUserById(id)
	if user.ID == 0 {
		return
	}

	guid := u.r.GetGuid(unionid)
	if guid == "" {
		guid = util2.UUID()
	}

	user.UnionId = unionid
	user.GUID = guid
	err := u.r.Save(&user)
	if err != nil {
		app.Logger.Error("更新unionid失败", id, unionid, err)
	}
}
func (u UserService) GetUserBy(by repository2.GetUserBy) (*entity.User, error) {
	user := repository2.DefaultUserRepository.GetUserBy(by)
	return &user, nil
}
func (u UserService) FindOrCreateByMobile(mobile string) (*entity.User, error) {
	user := repository2.DefaultUserRepository.GetUserBy(repository2.GetUserBy{
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
	user := repository2.DefaultUserRepository.GetUserById(userId)

	if user.ID == 0 || user.PhoneNumber == "" {
		return &entity.User{}, nil
	}

	sourceUer := repository2.DefaultUserRepository.GetUserBy(repository2.GetUserBy{
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
func (u UserService) BindPhoneByCode(userId int64, code string) error {
	userInfo := u.r.GetUserById(userId)
	if userInfo.ID == 0 {
		return errors.New("未查到用户信息")
	}

	phoneResult, err := app.Weapp.NewPhonenumber().GetPhoneNumber(&phonenumber.GetPhoneNumberRequest{
		Code: code,
	})
	if err != nil {
		return err
	}
	userInfo.PhoneNumber = phoneResult.Data.PhoneNumber
	return u.r.Save(&userInfo)
}
func (u UserService) BindPhoneByIV(param BindPhoneByIVParam) error {
	userInfo := u.r.GetUserById(param.UserId)
	if userInfo.ID == 0 {
		return errors.New("未查到用户信息")
	}

	session, err := DefaultSessionService.FindSessionByOpenId(userInfo.OpenId)

	if err != nil {
		return err
	}
	if session.ID == 0 {
		return errno.ErrAuth
	}

	decryptResult, err := app.Weapp.DecryptMobile(session.SessionKey, param.EncryptedData, param.IV)
	if err != nil {
		return err
	}
	userInfo.PhoneNumber = decryptResult.PhoneNumber
	return u.r.Save(&userInfo)
}
