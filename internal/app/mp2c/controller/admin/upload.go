package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/admin/admtypes"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultUploadController = UploadController{}

type UploadController struct {
}

func (UploadController) GetUploadTokenInfo(ctx *gin.Context) (gin.H, error) {
	form := admtypes.GetUploadTokenInfoForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthAdmin(ctx)

	info, err := service.DefaultUploadService.CreateUploadToken(int64(user.ID), 2, form.Scene)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"info": info,
	}, err
}
