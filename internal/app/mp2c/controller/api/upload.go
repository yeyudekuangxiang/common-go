package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/upload"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"net/http"
	"path"
	"strings"
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
			return nil, errno.ErrCommon.WithMessage("请选择文件")
		}
		return nil, err
	}

	if fileHeader.Size > 5*1024*1024 {
		return nil, errno.ErrCommon.WithMessage("文件大小不能超过5M")
	}

	fileExt := path.Ext(fileHeader.Filename)
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".jpeg" {
		return nil, errno.ErrCommon.WithMessage("仅支持png、jpg格式的图片")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	user := apiutil.GetAuthUser(ctx)
	imgUrl, err := upload.DefaultUploadService.UploadOcrImage(user.OpenId, file, fileHeader.Filename, form.PointCollectType)
	return gin.H{
		"imgUrl": imgUrl,
	}, err
}

func (UploadController) MultipartUploadImage(ctx *gin.Context) (gin.H, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, errno.ErrCommon.WithMessage("请选择文件")
		}
		return nil, err
	}
	//fileExt := path.Ext(fileHeader.Filename)
	/*if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".jpeg" {
		return nil, errors.New("仅支持png、jpg格式的图片")
	}*/

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	user := apiutil.GetAuthUser(ctx)
	imgUrl, err := upload.DefaultUploadService.MultipartUploadOcrImage(user.OpenId, file, fileHeader.Filename)
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

	info, err := upload.DefaultUploadService.CreateUploadToken(user.ID, 1, form.Scene)
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

	err := upload.DefaultUploadService.UploadCallback(srv_types.UploadCallbackParam{
		LogId:    logId,
		Filename: form.Filename,
		Size:     form.Size,
		MimeType: form.MimeType,
		Height:   form.Height,
		Width:    form.Width,
	})

	return nil, err
}
func (UploadController) UploadImage(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	if user.Auth != 1 {
		return nil, errno.ErrCommon.WithMessage("无权限")
	}
	userPlatform, _, err := service.DefaultUserService.FindOneUserPlatformByGuid(ctx.Request.Context(), user.GUID, entity.UserPlatformWxMiniApp)
	if err != nil {
		return nil, err
	}

	form := UploadImageForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, errno.ErrCommon.WithMessage("请选择文件")
		}
		return nil, err
	}
	uploadScene, err := upload.DefaultUploadSceneService.FindUploadScene(srv_types.FindSceneParam{
		Scene: strings.ToLower(form.ImageScene),
	})
	if err != nil || uploadScene.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("上传场景错误")
	}

	if fileHeader.Size > 5*1024*1024 {
		return nil, errno.ErrCommon.WithMessage("文件大小不能超过5M")
	}
	fileExt := path.Ext(fileHeader.Filename)
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".jpeg" {
		return nil, errno.ErrCommon.WithMessage("仅支持png、jpg格式的图片")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err = validator.CheckMediaWithOpenId(userPlatform.Openid, base64.StdEncoding.EncodeToString(data)); err != nil {
		return nil, err
	}

	imgUrl, err := upload.DefaultUploadService.UploadImage(userPlatform.Openid, file, fileHeader.Filename, uploadScene.OssDir)
	return gin.H{
		"imgUrl": imgUrl,
	}, err
}
func (UploadController) GetUploadSTSTokenInfo(ctx *gin.Context) (gin.H, error) {

	form := api_types.GetUploadTokenInfoForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)

	info, err := upload.DefaultUploadService.CreateStsToken(user.ID, 1, form.Scene)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"info": info,
	}, err
}
