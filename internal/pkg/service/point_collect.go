package service

import (
	"errors"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"strings"
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
	powerReplaceRuleOne = []string{
		"订单编号",
	}
	PointCollectBaiGuoYuanOne = []string{
		"购物袋",
		"袋",
	}
	PointCollectBaiGuoYuanTwo = []string{
		"百果园",
	}
)
var DefaultPointCollectService = PointCollectService{}

type PointCollectType string

const (
	PointCollectCoffeeCupType    PointCollectType = "COFFEE_CUP"
	PointCollectBikeRideType     PointCollectType = "BIKE_RIDE"
	PointCollectPowerReplaceType PointCollectType = "POWER_REPLACE"
	PointCollectDiDiType         PointCollectType = "DIDI"
	PointCollectReducePlastic    PointCollectType = "REDUCE_PLASTIC"
)

type PointCollectService struct {
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

func (srv PointCollectService) validatePowerReplaceImage(imageUrl string) (bool, []string, error) {
	results, err := DefaultOCRService.Scan(imageUrl)
	fmt.Println(results, err)
	if err != nil {
		return false, nil, err
	}
	result := srv.validatePointRule(results, powerReplaceRuleOne)
	if result != "" {
		return true, []string{result}, nil
	}
	return false, nil, nil
}

func (srv PointCollectService) validateReducePlasticImage(imageUrl string) (bool, []string, error) {
	results, err := DefaultOCRService.Scan(imageUrl)
	fmt.Println(results, err)
	if err != nil {
		return false, nil, err
	}
	result := srv.validatePointRule(results, PointCollectBaiGuoYuanOne)
	result2 := srv.validatePointRule(results, PointCollectBaiGuoYuanTwo)
	if result == "" && result2 != "" {
		return true, []string{result}, nil
	}
	return false, nil, nil
}

func (srv PointCollectService) CollectBikeRide(openId string, risk int, imageUrl string) (int, error) {
	err := DefaultOCRService.CheckIdempotent(openId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService.CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).
		CheckLimit(entity.POINT_BIKE_RIDE, openId)
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

	_, err = NewPointCollectHistoryService(context.NewMioContext()).CreateHistory(CreateHistoryParam{
		OpenId:          openId,
		TransactionType: entity.POINT_BIKE_RIDE,
		Info:            fmt.Sprintf("bikeRide=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加骑行更酷记录失败", openId, imageUrl, err)
	}

	value := entity.PointCollectValueMap[entity.POINT_BIKE_RIDE]
	_, err = NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openId,
		Type:         entity.POINT_BIKE_RIDE,
		BizId:        util.UUID(),
		ChangePoint:  int64(value),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})
	return value, err
}

func (srv PointCollectService) CollectCoffeeCup(openId string, risk int, imageUrl string) (int, error) {
	err := DefaultOCRService.CheckIdempotent(openId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService.CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).
		CheckLimit(entity.POINT_COFFEE_CUP, openId)
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

	_, err = NewPointCollectHistoryService(context.NewMioContext()).CreateHistory(CreateHistoryParam{
		OpenId:          openId,
		TransactionType: entity.POINT_COFFEE_CUP,
		Info:            fmt.Sprintf("coffeeCup=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加自带咖啡杯记录失败", openId, imageUrl, err)
	}

	value := entity.PointCollectValueMap[entity.POINT_COFFEE_CUP]
	_, err = NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openId,
		Type:         entity.POINT_COFFEE_CUP,
		BizId:        util.UUID(),
		ChangePoint:  int64(value),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})
	return value, err
}

func (srv PointCollectService) CollectPowerReplace(openId string, risk int, imageUrl string) (int, error) {
	err := DefaultOCRService.CheckIdempotent(openId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService.CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).
		CheckLimit(entity.POINT_POWER_REPLACE, openId)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("达到当日该类别最大积分限制")
	}

	valid, result, err := srv.validatePowerReplaceImage(imageUrl)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errors.New("不是有效的图片")
	}

	_, err = NewPointCollectHistoryService(context.NewMioContext()).CreateHistory(CreateHistoryParam{
		OpenId:          openId,
		TransactionType: entity.POINT_POWER_REPLACE,
		Info:            fmt.Sprintf("powerReplace=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加电车换电记录失败", openId, imageUrl, err)
	}

	value := entity.PointCollectValueMap[entity.POINT_POWER_REPLACE]
	_, err = NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openId,
		Type:         entity.POINT_POWER_REPLACE,
		BizId:        util.UUID(),
		ChangePoint:  int64(value),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})
	return value, err
}

func (srv PointCollectService) CollectReducePlastic(openId string, risk int, imageUrl string) (int, error) {

	err := DefaultOCRService.CheckIdempotent(openId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService.CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).
		CheckLimit(entity.POINT_REDUCE_PLASTIC, openId)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("达到当日该类别最大积分限制")
	}

	valid, result, err := srv.validateReducePlasticImage(imageUrl)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errors.New("不是有效的图片")
	}

	_, err = NewPointCollectHistoryService(context.NewMioContext()).CreateHistory(CreateHistoryParam{
		OpenId:          openId,
		TransactionType: entity.POINT_REDUCE_PLASTIC,
		Info:            fmt.Sprintf("reducePlastic=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加环保减塑记录失败", openId, imageUrl, err)
	}

	value := entity.PointCollectValueMap[entity.POINT_REDUCE_PLASTIC]
	_, err = NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openId,
		Type:         entity.POINT_REDUCE_PLASTIC,
		BizId:        util.UUID(),
		ChangePoint:  int64(value),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})
	return value, err
}
