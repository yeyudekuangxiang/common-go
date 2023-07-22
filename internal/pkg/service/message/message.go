package message

import (
	contextRedis "context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/medivhzhan/weapp/v3/request"
	"github.com/medivhzhan/weapp/v3/subscribemessage"
	"gitlab.miotech.com/miotech-application/backend/common-go/message"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/carbonpk/carbonpk"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/user/cmd/rpc/user"
	"golang.org/x/net/context"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/pkg/errno"

	"strconv"
	"time"
)

type MessageService struct {
}

// SendMiniSubMessage  小程序订阅消息发送
func (srv *MessageService) SendMiniSubMessage(toUser string, page string, template IMiniSubTemplate) (int, error) {
	zhuGeAttr := make(map[string]interface{}, 0) //诸葛打点
	zhuGeAttr["openid"] = toUser
	redisTemplateKey := fmt.Sprintf(config.RedisKey.MessageLimitByTemplate, template.TemplateId(), time.Now().Format("20060102"))
	if !util.DefaultLock.Lock(redisTemplateKey, time.Minute*5) {
		return 0, nil
	}
	defer util.DefaultLock.UnLock(redisTemplateKey)
	//redisUserKey := fmt.Sprintf(config.RedisKey.MessageLimitByUser, time.Now().Format("20060102"))
	defer track.DefaultZhuGeService().Track(config.ZhuGeEventName.MessageMiniSubscribe, toUser, zhuGeAttr)

	//发送限制
	templateSendCount := app.Redis.ZScore(contextRedis.Background(), redisTemplateKey, toUser).Val()
	//userSendCount := app.Redis.ZScore(contextRedis.Background(), redisUserKey, toUser).Val()
	if templateSendCount >= template.SendMixCount() {
		zhuGeAttr["错误码"] = -1
		zhuGeAttr["错误信息"] = "同一模板每人每天最多接收1条消息"
		return -1, errno.ErrCommon.WithMessage("同一模板每人每天最多接收1条消息")
	}

	/*if userSendCount >= 2 {
		zhuGeAttr["错误码"] = -2
		zhuGeAttr["错误信息"] = "每人每天最多收到2个不同类型模板消息"
		return -2, errors.New("每人每天最多收到2个不同类型模板消息")
	}*/

	var ret *request.CommonError
	err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		ret, err = app.Weapp.NewSubscribeMessage().Send(&subscribemessage.SendRequest{
			ToUser:           toUser,
			TemplateID:       template.TemplateId(),
			Page:             page,
			MiniprogramState: subscribemessage.MiniprogramStateFormal,
			Data:             template.ToData(),
		})
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(ret.ErrCode)
	}, 1)

	if err != nil {
		app.Logger.Infof("小程序订阅消息发送失败，http层，模版%s，toUser%s，错误信息%s", template.TemplateId(), toUser, err.Error())
		zhuGeAttr["错误码"] = -3
		zhuGeAttr["错误信息"] = err.Error()
		return -3, err
	}
	zhuGeAttr["错误码"] = ret.ErrCode
	zhuGeAttr["错误信息"] = ret.ErrMSG
	if ret.ErrCode != 43101 && ret.ErrCode != 0 {
		app.Logger.Infof("小程序订阅消息发送失败，业务层，模版%s，toUser%s，错误信息%s", template.TemplateId(), toUser, ret.GetResponseError().Error())
		return ret.ErrCode, ret.GetResponseError()
	}
	if ret.ErrCode == 0 {
		app.Redis.ZIncrBy(contextRedis.Background(), redisTemplateKey, 1, toUser) //同一模板每人每天最多接收1条消息
		//app.Redis.ZIncrBy(contextRedis.Background(), redisUserKey, 1, toUser).Err()     //每人每天最多收到2个不同类型模板消息
	}
	return ret.ErrCode, nil
}

//GetTemplateId 根据场景获取模版id
func (srv *MessageService) GetTemplateId(openid string, scene string) (templateIds []string) {
	var redisTemplateKey string
	switch scene {
	case "topic":
		templateIds = append(templateIds, config.MessageTemplateIds.TopicPass, config.MessageTemplateIds.TopicCarefullyChosen)
		redisTemplateKey = fmt.Sprintf(config.RedisKey.MessageLimitTopicShow, time.Now().Format("20060102"))
		break
	case "platform":
		templateIds = append(templateIds, config.MessageTemplateIds.ChangePoint, config.MessageTemplateIds.SignRemind, config.MessageTemplateIds.OrderDeliver)
		redisTemplateKey = fmt.Sprintf(config.RedisKey.MessageLimitPlatformShow, time.Now().Format("20060102"))
		srv.ExtensionSignTime(openid)
		break
	case "carbonpk":
		templateIds = append(templateIds, config.MessageTemplateIds.PunchClockRemind)
		redisTemplateKey = fmt.Sprintf(config.RedisKey.MessageLimitCarbonPkShow, time.Now().Format("20060102"))
		return
	case "quiz":
		templateIds = append(templateIds, config.MessageTemplateIds.QuizRemind)
		redisTemplateKey = fmt.Sprintf(config.RedisKey.QuizMessageRemind)
		return
	case "charge":
		templateIds = append(templateIds, message.MessageTemplateIds.ChargeOrder)
		redisTemplateKey = fmt.Sprintf(config.RedisKey.MessageLimitChargeRemind, time.Now().Format("20060102"))
	default:
		break
	}
	showCount := app.Redis.ZScore(contextRedis.Background(), redisTemplateKey, openid).Val()
	if config.Config.App.Env == "prod" {
		if showCount >= 1 && scene != "charge" {
			return []string{}
		}
	}
	if redisTemplateKey != "" {
		app.Redis.ZIncrBy(contextRedis.Background(), redisTemplateKey, 1, openid) //同一模板每人每天最多接收1条消息
	}
	return
}

//ExtensionSignTime 签到时间设置提醒时间
func (srv MessageService) ExtensionSignTime(openId string) {
	add := time.Now().Add(24 * time.Hour).Unix()
	app.Redis.ZRem(contextRedis.Background(), config.RedisKey.MessageSignUser, openId)
	app.Redis.ZAdd(contextRedis.Background(), config.RedisKey.MessageSignUser, &redis.Z{Score: float64(add), Member: openId})
}

//DelExtensionSignTime 删除用户签到提醒
func (srv MessageService) DelExtensionSignTime(openid string) {
	messageSignUser := config.RedisKey.MessageSignUser
	app.Redis.ZRem(contextRedis.Background(), messageSignUser, openid)
}

//SendMessageToSignUser 给用户发送签到提醒
func (srv MessageService) SendMessageToSignUser() {
	if !util.DefaultLock.Lock("SendMessageToSignUser", time.Minute*5) {
		return
	}
	defer util.DefaultLock.UnLock("SendMessageToSignUser")
	messageSignUserKey := config.RedisKey.MessageSignUser
	time := time.Now().Unix()
	op := &redis.ZRangeBy{
		Max:    strconv.FormatInt(time, 10),
		Min:    "0",
		Offset: 0,    //类似sql的limit, 表示开始偏移量
		Count:  5000, //默认一次跑5000条数据，因为是每10分钟跑一次，根据现在的日活，5000是可以的，后续可以增加或者进行分页处理
	}
	list, err := app.Redis.ZRevRangeByScore(contextRedis.Background(), messageSignUserKey, op).Result()
	if err != nil {
		return
	}
	message := MiniSignRemindTemplate{
		ActivityName: "每日签到",
		Tip:          "低碳打卡有惊喜，快来解锁今日福利吧",
	}
	successCount := 0
	failCount := 0
	refuseCount := 0
	for _, openid := range list {
		code, _ := srv.SendMiniSubMessage(openid, config.MessageJumpUrls.SignRemind, message)
		if code == 0 {
			//发送成功，提醒时间延长24小时
			srv.ExtensionSignTime(openid)
			successCount++
			continue
		} else if code == 43101 {
			//用户拒绝提醒，删除用户提醒
			srv.DelExtensionSignTime(openid) //删除提醒
			refuseCount++
			continue
		} else {
			srv.DelExtensionSignTime(openid) //删除提醒
			failCount++
			//其他错误情况，不做任何处理,会在下次定时器运行时，重新提醒
		}
	}
	app.Logger.Infof("签到消息发送，总条数%d,成功%d,失败%d,拒绝%d", len(list), successCount, failCount, refuseCount)
}

func (srv MessageService) SendMessageToCarbonPk() {
	if !util.DefaultLock.Lock("sendMessageToCarbonPk", time.Minute*5) {
		return
	}
	defer util.DefaultLock.UnLock("sendMessageToCarbonPk")
	redisKey := config.RedisKey.CarbonPkRemindUser
	list := app.Redis.SMembersMap(contextRedis.Background(), redisKey)
	template := MiniClockRemindTemplate{
		Title:   "低碳打卡挑战赛",
		Name:    "打卡提醒",
		Date:    "",
		Content: "快来完成今天打卡吧！惊喜离你越来越近了哦",
		Tip:     "低碳生活挑战，一起每天打卡吧",
	}

	for s := range list.Val() {
		uid, err := strconv.ParseInt(s, 10, 64)
		getUserById, err := app.RpcService.UserRpcSrv.FindUserByID(context.Background(), &user.FindUserByIDReq{
			UserId: uid,
		})
		if err != nil {
			app.Logger.Infof("低碳打卡挑战，小程序订阅消息发送失败，模版%s，uid%d，错误信息%s", template.TemplateId(), uid, err.Error())
			continue
		}
		if !getUserById.GetExist() {
			app.Logger.Infof("低碳打卡挑战，小程序订阅消息发送失败,用户不存在，模版%s，uid%d，错误信息%s", template.TemplateId(), uid, err.Error())
			continue
		}
		days, err := app.RpcService.CarbonPkRpcSrv.TotalCarbonPkDays(context.Background(), &carbonpk.TotalCarbonPkDaysReq{
			UserId: uid,
		})
		if err != nil {
			app.Logger.Infof("低碳打卡挑战，小程序订阅消息发送失败,获取打卡天数失败，模版%s，uid%d，错误信息%s", template.TemplateId(), uid, err.Error())
			continue
		}
		template.Date = strconv.FormatInt(days.Total, 10) + "天"
		openid := getUserById.UserInfo.Openid
		var ret *request.CommonError
		err = app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
			ret, err = app.Weapp.NewSubscribeMessage().Send(&subscribemessage.SendRequest{
				ToUser:           openid, //oy_BA5IGl1JgkJKbD14wq_-Yorqw
				TemplateID:       template.TemplateId(),
				Page:             "/pages/activity/punch/start/index", //页面跳转
				MiniprogramState: subscribemessage.MiniprogramStateFormal,
				Data:             template.ToData(),
			})
			if err != nil {
				return false, err
			}
			return app.Weapp.IsExpireAccessToken(ret.ErrCode)
		}, 1)
		if err != nil {
			app.Logger.Infof("小程序订阅消息发送失败，http层，模版%s，toUser%s，错误信息%s", template.TemplateId(), openid, err.Error())
		}
	}
}

func (srv MessageService) SendMessageToQuiz() {
	if !util.DefaultLock.Lock("sendMessageToQuiz", time.Minute*5) {
		return
	}
	defer util.DefaultLock.UnLock("sendMessageToQuiz")
	QuizRemindKey := config.RedisKey.QuizMessageRemind
	ctx := contextRedis.Background()
	list, err := app.Redis.SMembersMap(ctx, QuizRemindKey).Result()
	if err != nil {
		app.Logger.Infof("答题挑战提醒,小程序订阅消息发送失败,错误信息%s", err.Error())
		return
	}

	template := MiniQuizRemindTemplate{
		Title:   "答题挑战提醒",
		Name:    "每日一题",
		Content: "今天的每日一题还没有做哦！",
	}
	uIds := make([]int64, 0)
	for id, _ := range list {
		uid, _ := strconv.ParseInt(id, 10, 64)
		uIds = append(uIds, uid)
	}
	userList, err := app.RpcService.UserRpcSrv.GetUserList(ctx, &user.GetUserListReq{
		UserIds: uIds,
	})
	if err != nil {
		app.Logger.Infof("答题挑战提醒,小程序订阅消息发送失败,模版: %s，错误信息: %s", template.TemplateId(), err.Error())
		return
	}

	for _, usr := range userList.GetList() {
		code, err := srv.SendMiniSubMessage(usr.Openid, config.MessageJumpUrls.QuizRemind, template)
		if err != nil || code != 0 {
			app.Logger.Infof("答题挑战提醒,小程序订阅消息发送失败,模版: %s,错误信息: %s, code: %d", template.TemplateId(), err.Error(), code)
			app.Redis.ZRem(ctx, QuizRemindKey, usr.GetId())
		}
	}
}
