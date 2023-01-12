package business

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/business/businesstypes"
	"mio/internal/pkg/service/upload"
	"mio/internal/pkg/util/apiutil"
)

var DefaultUploadController = UploadController{}

type UploadController struct {
}

func (UploadController) GetUploadTokenInfo(ctx *gin.Context) (gin.H, error) {

	form := businesstypes.GetUploadTokenInfoForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthBusinessUser(ctx)

	info, err := upload.DefaultUploadService.CreateUploadToken(user.ID, 3, form.Scene)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"info": info,
	}, err
}
