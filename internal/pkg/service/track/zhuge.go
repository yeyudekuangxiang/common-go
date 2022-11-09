package track

import (
	"fmt"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/repository"
	service "mio/internal/pkg/service/userExtend"
	"mio/pkg/errno"

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

	Debug bool
}

func DefaultZhuGeService() *ZhuGeService {
	return NewZhuGeService(zhuge.NewClient(config.Config.Zhuge.AppKey, config.Config.Zhuge.AppSecret), config.Config.App.Env == "prod")
}
func NewZhuGeService(client *zhuge.Client, open bool) *ZhuGeService {
	return &ZhuGeService{client: client, Open: open}
	//return &ZhuGeService{client: client, Open: true, Debug: true}
}

// TrackPoint 积分打点
func (srv ZhuGeService) TrackPoint(point srv_types.TrackPoint) {
	if !srv.Open {
		app.Logger.Info("诸葛打点已关闭", point)
		return
	}
	ip := ""
	//没传递，用户扩展表查
	userExtend, exist, err := service.DefaultUserExtendService.GetUserExtend(repository.GetUserExtendBy{OpenId: point.OpenId})
	if err == nil && exist == true {
		ip = userExtend.Ip
	}
	err = srv.client.Track(types.Event{
		Dt:    "evt",
		Pl:    "js",
		Debug: 0,
		Ip:    ip,
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
		"备注":   point.AdditionInfo,
	})

	if err != nil {
		app.Logger.Errorf("积分打点失败 %+v %+v", err, point)
	}
}

// TrackPoints (新) 积分打点
func (srv ZhuGeService) TrackPoints(point srv_types.TrackPoints) {
	if !srv.Open {
		app.Logger.Info("诸葛打点已关闭", point)
		return
	}
	ip := ""
	//没传递，用户扩展表查
	userExtend, exist, err := service.DefaultUserExtendService.GetUserExtend(repository.GetUserExtendBy{OpenId: point.OpenId})
	if err == nil && exist == true {
		ip = userExtend.Ip
	}

	err = srv.client.Track(types.Event{
		Dt:    "evt",
		Pl:    "js",
		Debug: 0,
		Ip:    ip,
		Pr: types.EventJs{
			Ct:   time.Now().UnixMilli(),
			Eid:  "积分变动",
			Cuid: point.OpenId,
			Sid:  time.Now().UnixMilli(),
		},
	}, map[string]interface{}{
		"积分类型": point.PointType,
		"变动方式": util.Ternary(point.ChangeType == "dec", "积分消耗", "积分获取").String(),
		"变动数量": util.Ternary(point.ChangeType == "dec", -int(point.Value), int(point.Value)).Int(),
		"是否失败": util.Ternary(point.IsFail, "操作失败", "操作成功").String(),
		"失败原因": point.FailMessage,
	})

	if err != nil {
		app.Logger.Errorf("积分打点失败 %+v %+v", err, point)
	}
}

func (srv ZhuGeService) TrackBusinessPoints(point srv_types.TrackBusinessPoints) {
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
			Eid:  "企业版积分变动",
			Cuid: point.Uid,
			Sid:  time.Now().UnixMilli(),
		},
	}, map[string]interface{}{
		"用户编号": point.Uid,
		"变动类型": util.Ternary(point.ChangeType == "dec", "积分减少", "积分增加").String(),
		"变动数量": util.Ternary(point.ChangeType == "dec", -int(point.Value), int(point.Value)).Int(),
		"用户昵称": point.Nickname,
		"用户姓名": point.Username,
		"部门":   point.Department,
		"公司":   point.Company,
		"变动时间": point.ChangeTime.Format("2006/01/02 15:04:05"),
	})

	if err != nil {
		app.Logger.Errorf("企业版积分打点失败 %+v %+v", err, point)
	}
}
func (srv ZhuGeService) TrackBusinessCredit(credit srv_types.TrackBusinessCredit) {
	if !srv.Open {
		app.Logger.Info("诸葛打点已关闭", credit)
		return
	}
	err := srv.client.Track(types.Event{
		Dt:    "evt",
		Pl:    "js",
		Debug: 0,
		Pr: types.EventJs{
			Ct:   time.Now().UnixMilli(),
			Eid:  "企业版碳积分变动",
			Cuid: credit.Uid,
			Sid:  time.Now().UnixMilli(),
		},
	}, map[string]interface{}{
		"用户编号": credit.Uid,
		"变动类型": util.Ternary(credit.ChangeType == "dec", "积分减少", "积分增加").String(),
		"变动数量": util.Ternary(credit.ChangeType == "dec", -int(credit.Value), int(credit.Value)).Int(),
		"用户昵称": credit.Nickname,
		"用户姓名": credit.Username,
		"部门":   credit.Department,
		"公司":   credit.Company,
		"变动时间": credit.ChangeTime.Format("2006/01/02 15:04:05"),
	})

	if err != nil {
		app.Logger.Errorf("企业版碳积分打点失败 %+v %+v", err, credit)
	}
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
		Debug: util.Ternary(srv.Debug, int(1), int(0)).Int(),
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
		Debug: util.Ternary(srv.Debug, int(1), int(0)).Int(),
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
		app.Logger.Errorf("积分打点失败 %+v %v %+v", eventName, openId, attr)
	}
	return nil
}
