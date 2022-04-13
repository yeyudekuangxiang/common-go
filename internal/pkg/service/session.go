package service

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
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
