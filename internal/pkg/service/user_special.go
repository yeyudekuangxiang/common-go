package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultUserSpecialService = NewUserSpecialService(repository.DefaultUserSpecialRepository)

func NewUserSpecialService(r repository.UserSpecialRepository) UserSpecialService {
	return UserSpecialService{
		r: r,
	}
}

type UserSpecialService struct {
	r repository.UserSpecialRepository
}

// GetSpecialUserByPhone 查询用户信息
func (u UserSpecialService) GetSpecialUserByPhone(phone string) entity.UserSpecial {
	return u.r.GetUserSpecialByPhone(phone)
}

func (u UserSpecialService) Save(special *entity.UserSpecial) error {
	return u.r.Save(special)
}
