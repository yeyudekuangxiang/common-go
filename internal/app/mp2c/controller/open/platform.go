package open

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	platformUtil "mio/pkg/platform"
	"strings"
)

var DefaultPlatformController = PlatformController{}

type PlatformController struct {
}

func (receiver PlatformController) BindPlatformUser(ctx *gin.Context) (gin.H, error) {
	//接收参数 platformKey, phone
	form := platformForm{}
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

	if user.ID == 0 {
		return nil, errors.New("用户未登陆")
	}

	//保存渠道用户记录
	sceneUser := service.DefaultBdSceneUserService.FindPlatformUser(user.OpenId, form.PlatformKey)
	if sceneUser.ID == 0 {
		sceneUser.PlatformKey = scene.Ch
		sceneUser.PlatformUserId = form.MemberId
		sceneUser.Phone = user.PhoneNumber
		sceneUser.OpenId = user.OpenId
		sceneUser.UnionId = user.UnionId
		err := service.DefaultBdSceneUserService.Create(sceneUser)
		if err != nil {
			app.Logger.Errorf("create db_scene_user error:%s", err.Error())
		}
		return nil, nil
	}
	return nil, nil
}

func (receiver PlatformController) SyncPoint(ctx *gin.Context) (gin.H, error) {
	form := platformForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	dst := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &dst)
	if err != nil {
		return nil, err
	}

	//查询渠道号
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errors.New("渠道查询失败")
	}

	//check sign
	if err := platformUtil.CheckSign(dst); err != nil {
		app.Logger.Errorf("校验sign失败: %s", err.Error())
		return nil, err
	}

	//check user
	user, _ := service.DefaultUserService.GetUserBy(repository.GetUserBy{Mobile: form.Mobile, Source: entity.UserSourceMio})
	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	method := scene.Ch
	if form.Method != "" {
		method = strings.ToLower(method) + "_" + strings.ToLower(form.Method)
	}
	t := entity.PlatformMethodMap[method]

	value := entity.PointCollectValueMap[t]

	_, err = service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      user.OpenId,
		Type:        t,
		BizId:       util.UUID(),
		ChangePoint: int64(value),
		AdminId:     0,
		Note:        t.Text(),
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
