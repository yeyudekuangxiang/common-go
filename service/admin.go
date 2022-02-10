package service

import (
	"mio/internal/util"
	"mio/model/auth"
	"mio/model/entity"
	"mio/repository"
)

var DefaultAdminService IAdminService = NewAdminService(repository.DefaultAdminRepository)

type IAdminService interface {
	//根据管理员id获取管理员信息
	GetAdminById(int) (*entity.Admin, error)
	//根据token获取管理员
	GetAdminByToken(string) (*entity.Admin, error)
}

func NewAdminService(r repository.IAdminRepository) AdminService {
	return AdminService{
		r: r,
	}
}

type AdminService struct {
	r repository.IAdminRepository
}

func (a AdminService) GetAdminByToken(token string) (*entity.Admin, error) {
	if util.IsOldToken(token) {
		var authAdmin auth.OldAdmin
		err := util.ParseOldToken(token, &authAdmin)
		if err != nil {
			return nil, err
		}
		return a.r.GetAdminById(authAdmin.UserId.Int())
	}
	var authAdmin auth.Admin
	err := util.ParseToken(token, &authAdmin)
	if err != nil {
		return nil, err
	}
	return a.r.GetAdminById(authAdmin.ID)
}

func (a AdminService) GetAdminById(id int) (*entity.Admin, error) {
	return a.r.GetAdminById(id)
}
