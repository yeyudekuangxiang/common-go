package service

import (
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
)

var DefaultAdminService IAdminService = NewAdminService(repository.DefaultAdminRepository)

type IAdminService interface {
	// GetAdminById 根据管理员id获取管理员信息
	GetAdminById(int) (*entity.Admin, error)
	// GetAdminByToken 根据token获取管理员
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
	var authAdmin auth.Admin
	err := util.ParseToken(token, &authAdmin)
	if err != nil {
		return nil, err
	}
	admin := a.r.GetAdminById(authAdmin.ID)
	return &admin, nil
}

func (a AdminService) GetAdminById(id int) (*entity.Admin, error) {
	admin := a.r.GetAdminById(id)
	return &admin, nil
}
