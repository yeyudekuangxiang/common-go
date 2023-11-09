package admin

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultAdminController = AdminController{}

type AdminController struct {
}

func (AdminController) GetAdminList(ctx *gin.Context) (gin.H, error) {
	list, err := service.DefaultSystemAdminService.GetAdminList(repository.GetAdminListBy{
		DeletedAt: &sql.NullTime{},
	})

	return gin.H{
		"list": list,
	}, err
}
func (AdminController) Login(ctx *gin.Context) (gin.H, error) {
	var form AdminLoginForm
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	token, err := service.DefaultSystemAdminService.Login(form.Account, form.Password)
	return gin.H{
		"token": token,
	}, err
}
func (AdminController) GetLoginAdminInfo(ctx *gin.Context) (gin.H, error) {
	admin := apiutil.GetAuthAdmin(ctx)
	return gin.H{
		"admin": admin,
	}, nil
}
