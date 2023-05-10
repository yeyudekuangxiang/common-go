package track

import (
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/commontool"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/repository"
	service "mio/internal/pkg/service/userExtend"
	"mio/pkg/errno"

	"gitlab.miotech.com/miotech-application/backend/common-go/zhuge"
	"gitlab.miotech.com/miotech-application/backend/common-go/zhuge/types"
	"time"
)

type ZhuGeService struct {
	client *zhuge.Client
	//是否开启打点
	Open bool

	Debug bool
}

func DefaultZhuGeService() *ZhuGeService {
	return NewZhuGeService(zhuge.NewClient(config.Config.Zhuge.AppKey, config.Config.Zhuge.AppSecret, 0), config.Config.App.Env == "prod")
}
func NewZhuGeService(client *zhuge.Client, open bool) *ZhuGeService {
	return &ZhuGeService{client: client, Open: open}
	//return &ZhuGeService{client: client, Open: true, Debug: true}
}
func (srv ZhuGeService) Track(eventName, openId string, attr map[string]interface{}) {
	if !srv.Open {
		app.Logger.Infof("积分打点失败 %+v %v %+v", eventName, openId, attr)
		return
	}
	ip := ""
	ipAttr, ok := attr["ip"]
	if ok && ipAttr != "" {
		ip = fmt.Sprintf("%s", ipAttr)
	} else {
		//没传递，用户扩展表查
		userExtend, exist, err := service.DefaultUserExtendService.GetUserExtend(repository.GetUserExtendBy{OpenId: openId})
		if err == nil && exist == true {
			ip = userExtend.Ip
		}
	}

	err := srv.client.Track(types.Event{
		Dt:    "evt",
		Pl:    "js",
		Debug: commontool.Ternary(srv.Debug, int(1), int(0)).Int(),
		Ip:    ip,
		Pr: types.EventJs{
			Ct:   time.Now().UnixMilli(),
			Eid:  eventName,
			Cuid: openId,
			Sid:  time.Now().UnixMilli(),
		},
	}, attr)

	if err != nil {
		app.Logger.Errorf("积分打点失败 %+v %v %+v", eventName, openId, attr)
	}
}

func (srv ZhuGeService) TrackWithErr(eventName, openId string, attr map[string]interface{}) error {
	if !srv.Open {
		app.Logger.Infof("积分打点失败 %+v %v %+v", eventName, openId, attr)
		return errno.ErrCommon.WithMessage("未开启打点，请联系开发")
	}
	ip := ""
	ipAttr, ok := attr["ip"]
	if ok && ipAttr != "" {
		ip = fmt.Sprintf("%s", ipAttr)
	} else {
		//没传递，用户扩展表查
		userExtend, exist, err := service.DefaultUserExtendService.GetUserExtend(repository.GetUserExtendBy{OpenId: openId})
		if err == nil && exist == true {
			ip = userExtend.Ip
		}
	}

	err := srv.client.Track(types.Event{
		Dt:    "evt",
		Pl:    "js",
		Debug: commontool.Ternary(srv.Debug, int(1), int(0)).Int(),
		Ip:    ip,
		Pr: types.EventJs{
			Ct:   time.Now().UnixMilli(),
			Eid:  eventName,
			Cuid: openId,
			Sid:  time.Now().UnixMilli(),
		},
	}, attr)

	if err != nil {
		return err
	}
	return nil
}
