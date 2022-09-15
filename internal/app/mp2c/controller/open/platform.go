package open

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultPlatformController = PlatformController{}

type PlatformController struct {
}

func (receiver PlatformController) BindPlatformUser(ctx *gin.Context) (gin.H, error) {
	//接收参数 platformKey, phone
	form := platform{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	//查询渠道号
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errors.New("渠道查询失败")
	}

	//查询手机号
	if user.ID == 0 {
		return nil, errors.New("用户未授权登陆")
	}

	//保存渠道用户
	sceneUser := service.DefaultBdSceneUserService.FindByCh(form.PlatformKey)
	if sceneUser.ID == 0 {
		sceneUser.Ch = scene.Ch
		sceneUser.SceneUserId = form.MemberId
		sceneUser.Phone = user.PhoneNumber
		sceneUser.OpenId = user.OpenId
		sceneUser.UnionId = user.UnionId
		err := service.DefaultBdSceneUserService.Create(sceneUser)
		app.Logger.Errorf("create db_scene_user error:%s", err.Error())
		return nil, nil
	}
	//调用第三方回调
	return nil, nil
}
