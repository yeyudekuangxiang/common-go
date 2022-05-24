package service

import (
	"errors"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
	"strings"
	"time"
)

var (
	coffeeCupRuleOne = []string{
		"自带杯",
	}
	coffeeCupRuleTwo = []string{
		"订单",
		"单号",
	}
	bikeRideRuleOne = []string{
		"骑行",
		"单车",
		"骑车",
		"bike",
		"出行",
		"哈啰",
		"摩拜",
		"青桔",
	}
)
var DefaultPointCollectService = PointCollectService{}

type PointCollectService struct {
}

func (srv PointCollectService) CollectCoffeeCup(openId string, imageUrl string) (int, error) {
	if !util.DefaultLock.Lock(fmt.Sprintf("CollectCoffeeCup%s", openId), time.Second*5) {
		return 0, errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(fmt.Sprintf("CollectCoffeeCup%s", openId))

	ok, err := DefaultPointTransactionCountLimitService.CheckLimit(entity.POINT_COFFEE_CUP, openId)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("达到当日该类别最大积分限制")
	}

	valid, result, err := srv.validateCoffeeCupImage(imageUrl)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errors.New("不是有效的自带杯图片")
	}

	_, err = DefaultPointCollectHistoryService.CreateHistory(CreateHistoryParam{
		OpenId:          openId,
		TransactionType: entity.POINT_COFFEE_CUP,
		Info:            fmt.Sprintf("coffeeCup=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加自带咖啡杯记录失败", openId, imageUrl, err)
	}

	value := entity.PointCollectValueMap[entity.POINT_COFFEE_CUP]
	_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:       openId,
		Type:         entity.POINT_COFFEE_CUP,
		Value:        value,
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})
	return value, err
}
func (srv PointCollectService) validateCoffeeCupImage(imageUrl string) (bool, []string, error) {
	results, err := DefaultOCRService.Scan(imageUrl)
	if err != nil {
		return false, nil, err
	}
	result1 := srv.validatePointRule(results, coffeeCupRuleOne)
	result2 := srv.validatePointRule(results, coffeeCupRuleTwo)

	if result1 != "" && result2 != "" {
		return true, []string{result1, result2}, nil
	}
	return false, nil, nil
}
func (srv PointCollectService) validatePointRule(texts []string, rules []string) string {
	for _, text := range texts {
		for _, rule := range rules {
			if strings.Contains(strings.ToLower(text), strings.ToLower(rule)) {
				return text
			}
		}
	}
	return ""
}
func (srv PointCollectService) validateBikeRideImage(imageUrl string) (bool, []string, error) {
	results, err := DefaultOCRService.Scan(imageUrl)
	if err != nil {
		return false, nil, err
	}
	result := srv.validatePointRule(results, bikeRideRuleOne)
	if result != "" {
		return true, []string{result}, nil
	}
	return false, nil, nil
}
func (srv PointCollectService) CollectBikeRide(openId string, imageUrl string) (int, error) {
	if !util.DefaultLock.Lock(fmt.Sprintf("CollectBikeRide%s", openId), time.Second*5) {
		return 0, errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(fmt.Sprintf("CollectBikeRide%s", openId))

	ok, err := DefaultPointTransactionCountLimitService.CheckLimit(entity.POINT_BIKE_RIDE, openId)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("达到当日该类别最大积分限制")
	}

	valid, result, err := srv.validateBikeRideImage(imageUrl)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errors.New("不是有效的单车图片")
	}

	_, err = DefaultPointCollectHistoryService.CreateHistory(CreateHistoryParam{
		OpenId:          openId,
		TransactionType: entity.POINT_BIKE_RIDE,
		Info:            fmt.Sprintf("bikeRide=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加骑行更酷记录失败", openId, imageUrl, err)
	}

	value := entity.PointCollectValueMap[entity.POINT_BIKE_RIDE]
	_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:       openId,
		Type:         entity.POINT_BIKE_RIDE,
		Value:        value,
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})
	return value, err
}
