package message

import (
	contextRedis "context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/medivhzhan/weapp/v3/subscribemessage"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/track"
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
	redisUserKey := fmt.Sprintf(config.RedisKey.MessageLimitByUser, time.Now().Format("20060102"))

	defer track.DefaultZhuGeService().Track(config.ZhuGeEventName.MessageMiniSubscribe, toUser, zhuGeAttr)

	//发送限制

	templateSendCount := app.Redis.ZScore(contextRedis.Background(), redisTemplateKey, toUser).Val()
	userSendCount := app.Redis.ZScore(contextRedis.Background(), redisUserKey, toUser).Val()

	if templateSendCount >= 1 {
		zhuGeAttr["错误码"] = -1
		zhuGeAttr["错误信息"] = "同一模板每人每天最多接收1条消息"
		return -1, errors.New("同一模板每人每天最多接收1条消息")
	}

	if userSendCount >= 2 {
		zhuGeAttr["错误码"] = -2
		zhuGeAttr["错误信息"] = "每人每天最多收到2个不同类型模板消息"
		return -2, errors.New("每人每天最多收到2个不同类型模板消息")
	}

	ret, err := app.Weapp.NewSubscribeMessage().Send(&subscribemessage.SendRequest{
		ToUser:           toUser,
		TemplateID:       template.TemplateId(),
		Page:             page,
		MiniprogramState: subscribemessage.MiniprogramStateFormal,
		Data:             template.ToData(),
	})
	if err != nil {
		app.Logger.Info("小程序订阅消息发送失败，http层，模版%s，toUser%s，错误信息%s", template.TemplateId(), toUser, err.Error())
		zhuGeAttr["错误码"] = -3
		zhuGeAttr["错误信息"] = err.Error()
		return -3, err
	}
	zhuGeAttr["错误码"] = ret.ErrCode
	zhuGeAttr["错误信息"] = ret.ErrMSG
	if ret.ErrCode != 43101 && ret.ErrCode != 0 {
		app.Logger.Info("小程序订阅消息发送失败，业务层，模版%s，toUser%s，错误信息%s", template.TemplateId(), toUser, ret.GetResponseError().Error())
		return ret.ErrCode, ret.GetResponseError()
	}
	if ret.ErrCode == 0 {
		app.Redis.ZIncrBy(contextRedis.Background(), redisTemplateKey, 1, toUser).Err() //同一模板每人每天最多接收1条消息
		app.Redis.ZIncrBy(contextRedis.Background(), redisUserKey, 1, toUser).Err()     //每人每天最多收到2个不同类型模板消息
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
	default:
		break
	}
	showCount := app.Redis.ZScore(contextRedis.Background(), redisTemplateKey, openid).Val()
	if showCount >= 1 {
		return []string{}
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
		code, messageErr := srv.SendMiniSubMessage(openid, config.MessageJumpUrls.SignRemind, message)
		if messageErr == nil && code == 0 {
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
			failCount++
			//其他错误情况，不做任何处理,会在下次定时器运行时，重新提醒
		}
	}
	app.Logger.Infof("签到消息发送，总条数%d,成功%d,失败%d,拒绝%d", len(list), successCount, failCount, refuseCount)
}

func (srv MessageService) SendMessageToSignUserTest() (gin.H, error) {
	messageSignUserKey := config.RedisKey.MessageSignUser
	time := time.Now().Unix()
	op := &redis.ZRangeBy{
		Max:    strconv.FormatInt(time, 10),
		Min:    "0",
		Offset: 0,    //类似sql的limit, 表示开始偏移量
		Count:  5000, //默认一次跑5000条数据，因为是每10分钟跑一次，根据现在的日活，5000是可以的，后续可以增加或者进行分页处理
	}
	list, err := app.Redis.ZRevRangeByScore(contextRedis.Background(), messageSignUserKey, op).Result()

	return gin.H{
		"code": list,
		"err":  err,
	}, nil

	message := MiniSignRemindTemplate{
		ActivityName: "每日签到",
		Tip:          "低碳打卡有惊喜，快来解锁今日福利吧",
	}
	successCount := 0
	failCount := 0
	refuseCount := 0
	for _, openid := range list {
		code, messageErr := srv.SendMiniSubMessage(openid, config.MessageJumpUrls.SignRemind, message)
		if messageErr == nil && code == 0 {
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
			failCount++
			//其他错误情况，不做任何处理,会在下次定时器运行时，重新提醒
		}
	}
	app.Logger.Infof("签到消息发送，总条数%d,成功%d,失败%d,拒绝%d", len(list), successCount, failCount, refuseCount)
	return gin.H{}, nil
}
