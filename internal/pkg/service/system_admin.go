package service

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
)

var DefaultSystemAdminService ISystemAdminService = NewSystemAdminService(repository.DefaultSystemAdminRepository)

type ISystemAdminService interface {
	// GetAdminById 根据管理员id获取管理员信息
	GetAdminById(int) (*entity.SystemAdmin, error)
	// GetAdminByToken 根据token获取管理员
	GetAdminByToken(string) (*entity.SystemAdmin, bool, error)
	GetAdminList(by repository.GetAdminListBy) ([]entity.SystemAdmin, error)
	Login(account, password string) (string, error)
}

func NewSystemAdminService(r repository.ISystemAdminRepository) SystemAdminService {
	return SystemAdminService{
		r: r,
	}
}

type SystemAdminService struct {
	r repository.ISystemAdminRepository
}

func (a SystemAdminService) GetAdminByToken(token string) (*entity.SystemAdmin, bool, error) {
	var authAdmin auth.Admin
	err := util.ParseToken(token, &authAdmin)
	if err != nil {
		return nil, false, err
	}
	admin := a.r.GetAdminById(authAdmin.MioAdminID)
	if admin.ID == 0 {
		return nil, false, nil
	}
	return &admin, true, nil
}

func (a SystemAdminService) GetAdminById(id int) (*entity.SystemAdmin, error) {
	if id == 0 {
		return &entity.SystemAdmin{}, nil
	}
	admin := a.r.GetAdminById(id)
	return &admin, nil
}
func (a SystemAdminService) GetAdminList(by repository.GetAdminListBy) ([]entity.SystemAdmin, error) {
	return a.r.GetAdminList(by), nil
}
func (a SystemAdminService) Login(account, password string) (string, error) {
	admin := a.r.FindAdminBy(repository.FindAdminBy{
		Account: account,
	})
	if admin.ID == 0 {
		return "", errno.ErrAdminNotFound
	}

	if admin.Password != encrypttool.Md5(password) {
		return "", errno.ErrValidation
	}

	return util.CreateToken(auth.Admin{
		MioAdminID: admin.ID,
	})
}
