package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"net/http"
	"path"
)

var DefaultUploadController = UploadController{}

type UploadController struct {
}

func (UploadController) UploadPointCollectImage(ctx *gin.Context) (gin.H, error) {
	form := UploadPointCollectImageForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, errors.New("请选择文件")
		}
		return nil, err
	}

	if fileHeader.Size > 5*1024*1024 {
		return nil, errors.New("文件大小不能超过5M")
	}

	fileExt := path.Ext(fileHeader.Filename)
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".jpeg" {
		return nil, errors.New("仅支持png、jpg格式的图片")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	user := apiutil.GetAuthUser(ctx)
	imgUrl, err := service.DefaultUploadService.UploadOcrImage(user.OpenId, file, fileHeader.Filename, service.PointCollectType(form.PointCollectType))
	return gin.H{
		"imgUrl": imgUrl,
	}, err
}
func (UploadController) GetUploadTokenInfo(ctx *gin.Context) (gin.H, error) {

	form := api_types.GetUploadTokenInfoForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)

	info, err := service.DefaultUploadService.CreateUploadToken(user.ID, 1, form.Scene)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"info": info,
	}, err
}
func (UploadController) UploadCallback(ctx *gin.Context) (gin.H, error) {
	form := api_types.OssUploadCallbackForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	logId := ctx.Query("logId")
	if logId == "" {
		return nil, errno.ErrBind.WithCaller()
	}

	err := service.DefaultUploadService.UploadCallback(srv_types.UploadCallbackParam{
		LogId:    logId,
		Filename: form.Filename,
		Size:     form.Size,
		MimeType: form.MimeType,
		Height:   form.Height,
		Width:    form.Width,
	})

	return nil, err
}
