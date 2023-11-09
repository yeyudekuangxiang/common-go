package service

import (
	"database/sql"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"time"
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
	admin := a.r.GetAdminById(int(authAdmin.MioAdminID))
	if admin.ID == 0 {
		return nil, false, nil
	}
	if admin.Status != 1 {
		return nil, false, errno.ErrAuth.WithMessage("账号已被禁用")
	}
	if admin.DeletedAt.Valid {
		return nil, false, errno.ErrAuth.WithMessage("账号已被删除")
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
		Account:   account,
		DeletedAt: &sql.NullTime{},
	})
	if admin.ID == 0 {
		return "", errno.ErrAdminNotFound
	}
	if admin.Status != 1 {
		return "", errno.ErrCommon.WithMessage("账号已被禁用")
	}

	if admin.Password != encrypttool.Md5(password) {
		return "", errno.ErrValidation
	}

	return util.CreateToken(auth.Admin{
		Type:       "Admin",
		ID:         int64(admin.ID),
		MioAdminID: int64(admin.ID),
		CreatedAt:  time.Now().UnixMilli(),
	})
}
