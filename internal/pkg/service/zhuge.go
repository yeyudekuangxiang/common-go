package service

import (
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/zhuge"
	"mio/pkg/zhuge/types"
	"time"
)

type ZhuGeService struct {
	client *zhuge.Client
	//是否开启打点
	Open bool
}

func DefaultZhuGeService() *ZhuGeService {
	return NewZhuGeService(zhuge.NewClient(config.Config.Zhuge.AppKey, config.Config.Zhuge.AppSecret), config.Config.App.Env == "prod")
}
func NewZhuGeService(client *zhuge.Client, open bool) *ZhuGeService {
	return &ZhuGeService{client: client, Open: open}
}

// TrackPoint 积分打点
func (srv ZhuGeService) TrackPoint(point srv_types.TrackPoint) {
	if !srv.Open {
		app.Logger.Info("诸葛打点已关闭", point)
		return
	}
	err := srv.client.Track(types.Event{
		Dt:    "evt",
		Pl:    "js",
		Debug: 0,
		Pr: types.EventJs{
			Ct:   time.Now().UnixMilli(),
			Eid:  "积分变动",
			Cuid: point.OpenId,
			Sid:  time.Now().UnixMilli(),
		},
	}, map[string]interface{}{
		"积分类型": point.PointType.RealText(),
		"变动方式": util.Ternary(point.ChangeType == "dec", "积分消耗", "积分获取").String(),
		"变动数量": util.Ternary(point.ChangeType == "dec", -int(point.Value), int(point.Value)).Int(),
		"是否失败": util.Ternary(point.IsFail, "操作失败", "操作成功").String(),
		"失败原因": point.FailMessage,
	})

	if err != nil {
		app.Logger.Errorf("积分打点失败 %+v %+v", err, point)
	}
}
