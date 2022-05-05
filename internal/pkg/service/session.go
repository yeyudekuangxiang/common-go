package service

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/pkg/errno"
)

var DefaultSessionService = SessionService{}

type SessionService struct {
}

func (srv SessionService) FindSessionByOpenId(openId string) (*entity.Session, error) {
	session := entity.Session{}
	err := app.DB.Where("openid = ?", openId).Order("time desc").First(&session).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &session, nil
}

// MustGetSessionKey 获取失败时会返回 需要重新登陆错误码
func (srv SessionService) MustGetSessionKey(openid string) (string, error) {
	session, err := srv.FindSessionByOpenId(openid)
	if err != nil {
		return "", err
	}
	if session.ID == 0 {
		return "", errno.ErrAuth
	}

	return session.WechatSessionKey, nil
}
