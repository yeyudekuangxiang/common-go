package service

import (
	"fmt"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	carbonProducer "mio/internal/pkg/queue/producer/carbon"
	"mio/internal/pkg/queue/producer/growth_system"
	carbonmsg "mio/internal/pkg/queue/types/message/carbon"
	"mio/internal/pkg/queue/types/message/growthsystemmsg"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strconv"
	"strings"
)

var (
	coffeeCupRuleOne = []string{
		"自带杯",
		"cup discount",
		"cup",
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
	GreenTakeoutOne = []string{
		"无需餐具",
	}
	SustainablePackageDfOne = []string{
		"德芙",
		"Dove",
	}
)
var DefaultPointCollectService = PointCollectService{}

type PointCollectType string

const (
	PointCollectCoffeeCupType      PointCollectType = "COFFEE_CUP"
	PointCollectBikeRideType       PointCollectType = "BIKE_RIDE"
	PointCollectPowerReplaceType   PointCollectType = "POWER_REPLACE"
	PointCollectDiDiType           PointCollectType = "DIDI"
	PointCollectReducePlastic      PointCollectType = "REDUCE_PLASTIC"
	PointCollectGreenTakeout       PointCollectType = "GREEN_TAKE_OUT"
	PointCollectSustainablePackage PointCollectType = "SUSTAINABLE_PACKAGE"
)

type PointCollectService struct {
}

func (srv PointCollectService) validateCoffeeCupImage(imageUrl string) (bool, []string, error) {

	ocrSrv := DefaultOCRService()
	imageHash, err := ocrSrv.CheckImageScanCount(imageUrl, 1)
	if err != nil {
		return false, nil, err
	}
	results, err := ocrSrv.ScanWithHash(imageUrl, imageHash)

	if err != nil {
		return false, nil, err
	}
	result1 := srv.validatePointRule(results, coffeeCupRuleOne)
	//result2 := srv.validatePointRule(results, coffeeCupRuleTwo)

	if result1 != "" { //&& result2 != ""
		return true, []string{result1}, nil //result2
	}
	return false, nil, nil
}

func (srv PointCollectService) validateImage(imageUrl string, rules []string) (bool, []string, error) {
	ocrSrv := DefaultOCRService()
	imageHash, err := ocrSrv.CheckImageScanCount(imageUrl, 1)
	if err != nil {
		return false, nil, err
	}
	results, err := ocrSrv.ScanWithHash(imageUrl, imageHash)

	if err != nil {
		return false, nil, err
	}
	result1 := srv.validatePointRule(results, rules)
	if result1 != "" {
		return true, []string{result1}, nil
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
	ocrSrv := DefaultOCRService()
	imageHash, err := ocrSrv.CheckImageScanCount(imageUrl, 1)
	if err != nil {
		return false, nil, err
	}
	results, err := ocrSrv.ScanWithHash(imageUrl, imageHash)
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
	ocrSrv := DefaultOCRService()
	imageHash, err := ocrSrv.CheckImageScanCount(imageUrl, 1)
	if err != nil {
		return false, nil, err
	}
	results, err := ocrSrv.ScanWithHash(imageUrl, imageHash)
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
	ocrSrv := DefaultOCRService()
	imageHash, err := ocrSrv.CheckImageScanCount(imageUrl, 1)
	if err != nil {
		return false, nil, err
	}
	results, err := ocrSrv.ScanWithHash(imageUrl, imageHash)
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

func (srv PointCollectService) CollectBikeRide(uInfo entity.User, risk int, imageUrl string) (int, error) {
	err := DefaultOCRService().CheckIdempotent(uInfo.OpenId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService().CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).
		CheckLimit(entity.POINT_BIKE_RIDE, uInfo.OpenId)
	if err != nil {
		return 0, err
	}

	valid, result, err := srv.validateBikeRideImage(imageUrl)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errno.ErrCommon.WithMessage("不是有效的单车图片")
	}
	bizId := util.UUID()
	ctx := context.NewMioContext()
	var point int
	//减碳量
	_, err = NewCarbonTransactionService(ctx).Create(api_types.CreateCarbonTransactionDto{
		OpenId: uInfo.OpenId,
		Type:   entity.CARBON_BIKE_RIDE,
		Value:  1,
		Info:   fmt.Sprintf("%s", result),
		BizId:  bizId,
	})
	if err != nil {
		app.Logger.Error("添加骑行更酷减碳量失败", uInfo.OpenId, imageUrl, err)
	}

	//成长体系
	growth_system.GrowthSystemRide(growthsystemmsg.GrowthSystemParam{
		TaskType:    string(entity.POINT_BIKE_RIDE),
		TaskSubType: string(entity.POINT_BIKE_RIDE),
		UserId:      strconv.FormatInt(uInfo.ID, 10),
		TaskValue:   1,
	})

	if !ok {
		return point, errno.ErrCommon.WithMessage("今日次数以达到上限")
	}

	_, err = NewPointCollectHistoryService(ctx).CreateHistory(CreateHistoryParam{
		OpenId:          uInfo.OpenId,
		TransactionType: entity.POINT_BIKE_RIDE,
		Info:            fmt.Sprintf("bikeRide=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加骑行更酷记录失败", uInfo.OpenId, imageUrl, err)
	}

	point = entity.PointCollectValueMap[entity.POINT_BIKE_RIDE]
	_, err = NewPointService(ctx).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       uInfo.OpenId,
		Type:         entity.POINT_BIKE_RIDE,
		BizId:        bizId,
		ChangePoint:  int64(point),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})

	return point, err
}

func (srv PointCollectService) CollectCoffeeCup(uInfo entity.User, risk int, imageUrl string) (int, error) {
	err := DefaultOCRService().CheckIdempotent(uInfo.OpenId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService().CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).CheckLimit(entity.POINT_COFFEE_CUP, uInfo.OpenId)
	if err != nil {
		return 0, err
	}

	valid, result, err := srv.validateCoffeeCupImage(imageUrl)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errno.ErrCommon.WithMessage("不是有效的自带杯图片")
	}
	ctx := context.NewMioContext()
	bizId := util.UUID()
	var point int
	//减碳量
	_, err = NewCarbonTransactionService(ctx).Create(api_types.CreateCarbonTransactionDto{
		OpenId: uInfo.OpenId,
		Type:   entity.CARBON_COFFEE_CUP,
		Value:  1,
		Info:   fmt.Sprintf("%s", result),
		BizId:  bizId,
	})
	if err != nil {
		app.Logger.Error("添加自带杯减碳量失败", uInfo.OpenId, imageUrl, err)
	}
	//成长体系
	growth_system.GrowthSystemEPCoffeeCup(growthsystemmsg.GrowthSystemParam{
		TaskType:    string(entity.CARBON_COFFEE_CUP),
		TaskSubType: string(entity.CARBON_COFFEE_CUP),
		UserId:      strconv.FormatInt(uInfo.ID, 10),
		TaskValue:   1,
	})
	//每日上限检查
	if !ok {
		return point, errno.ErrCommon.WithMessage("今日次数已达到上限")
	}
	_, err = NewPointCollectHistoryService(ctx).CreateHistory(CreateHistoryParam{
		OpenId:          uInfo.OpenId,
		TransactionType: entity.POINT_COFFEE_CUP,
		Info:            fmt.Sprintf("coffeeCup=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加自带咖啡杯记录失败", uInfo.OpenId, imageUrl, err)
	}

	point = entity.PointCollectValueMap[entity.POINT_COFFEE_CUP]
	_, err = NewPointService(ctx).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       uInfo.OpenId,
		Type:         entity.POINT_COFFEE_CUP,
		BizId:        bizId,
		ChangePoint:  int64(point),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})

	return point, err
}

//绿色外卖

func (srv PointCollectService) CollectGreenTakeoutOne(uInfo entity.User, risk int, imageUrl string) (int, error) {
	err := DefaultOCRService().CheckIdempotent(uInfo.OpenId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService().CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).CheckLimit(entity.POINT_GREEN_TAKE_OUT, uInfo.OpenId)
	if err != nil {
		return 0, err
	}

	valid, result, err := srv.validateImage(imageUrl, GreenTakeoutOne)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errno.ErrCommon.WithMessage("不是有效的绿色外卖图片")
	}
	ctx := context.NewMioContext()
	bizId := util.UUID()
	var point int
	//减碳量
	_, err = NewCarbonTransactionService(ctx).Create(api_types.CreateCarbonTransactionDto{
		OpenId: uInfo.OpenId,
		Type:   entity.CARBON_GREEN_TAKE_OUT,
		Value:  1,
		Info:   fmt.Sprintf("%s", result),
		BizId:  bizId,
	})
	if err != nil {
		app.Logger.Error("添加绿色外卖减碳量失败", uInfo.OpenId, imageUrl, err)
	}
	//每日上限检查
	if !ok {
		return point, errno.ErrCommon.WithMessage("今日次数已达到上限")
	}
	_, err = NewPointCollectHistoryService(ctx).CreateHistory(CreateHistoryParam{
		OpenId:          uInfo.OpenId,
		TransactionType: entity.POINT_GREEN_TAKE_OUT,
		Info:            fmt.Sprintf("coffeeCup=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加绿色外卖杯记录失败", uInfo.OpenId, imageUrl, err)
	}

	point = entity.PointCollectValueMap[entity.POINT_GREEN_TAKE_OUT]
	_, err = NewPointService(ctx).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       uInfo.OpenId,
		Type:         entity.POINT_GREEN_TAKE_OUT,
		BizId:        bizId,
		ChangePoint:  int64(point),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})
	return point, err
}

//可持续包装

func (srv PointCollectService) CollectSustainablePackageOne(uInfo entity.User, risk int, imageUrl string) (int, error) {
	err := DefaultOCRService().CheckIdempotent(uInfo.OpenId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService().CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).CheckLimit(entity.POINT_SUSTAINABLE_PACKAGE, uInfo.OpenId)
	if err != nil {
		return 0, err
	}

	valid, result, err := srv.validateImage(imageUrl, SustainablePackageDfOne)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errno.ErrCommon.WithMessage("不是有效的可持续包装图片")
	}
	ctx := context.NewMioContext()
	bizId := util.UUID()
	var point int

	//每日上限检查
	if !ok {
		return point, errno.ErrCommon.WithMessage("今日次数已达到上限")
	}
	_, err = NewPointCollectHistoryService(ctx).CreateHistory(CreateHistoryParam{
		OpenId:          uInfo.OpenId,
		TransactionType: entity.POINT_SUSTAINABLE_PACKAGE,
		Info:            fmt.Sprintf("coffeeCup=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加可持续包装记录失败", uInfo.OpenId, imageUrl, err)
	}

	point = entity.PointCollectValueMap[entity.POINT_SUSTAINABLE_PACKAGE]
	_, err = NewPointService(ctx).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       uInfo.OpenId,
		Type:         entity.POINT_SUSTAINABLE_PACKAGE,
		BizId:        bizId,
		ChangePoint:  int64(point),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})

	//投递mq
	if err = carbonProducer.ChangeSuccessToQueue(carbonmsg.CarbonChangeSuccess{
		Openid:        uInfo.OpenId,
		UserId:        uInfo.ID,
		TransactionId: bizId,
		Type:          string(entity.POINT_SUSTAINABLE_PACKAGE),
		City:          "",
		Value:         0,
		Info:          fmt.Sprintf("%s", result),
	}); err != nil {
		app.Logger.Errorf("ChangeSuccessToQueue 投递失败:%v", err)
	}

	return point, err
}

// Deprecated: 使用ImageCollect代替
func (srv PointCollectService) CollectPowerReplace(uInfo entity.User, risk int, imageUrl string) (int, error) {
	err := DefaultOCRService().CheckIdempotent(uInfo.OpenId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService().CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).
		CheckLimit(entity.POINT_POWER_REPLACE, uInfo.OpenId)
	if err != nil {
		return 0, err
	}
	var point int
	valid, result, err := srv.validatePowerReplaceImage(imageUrl)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errno.ErrCommon.WithMessage("不是有效的图片")
	}

	if !ok {
		return point, errno.ErrCommon.WithMessage("今日次数以达到上限")
	}

	_, err = NewPointCollectHistoryService(context.NewMioContext()).CreateHistory(CreateHistoryParam{
		OpenId:          uInfo.OpenId,
		TransactionType: entity.POINT_POWER_REPLACE,
		Info:            fmt.Sprintf("powerReplace=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加电车换电记录失败", uInfo.OpenId, imageUrl, err)
	}

	point = entity.PointCollectValueMap[entity.POINT_POWER_REPLACE]
	_, err = NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       uInfo.OpenId,
		Type:         entity.POINT_POWER_REPLACE,
		BizId:        util.UUID(),
		ChangePoint:  int64(point),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})

	return point, err
}

func (srv PointCollectService) CollectReducePlastic(uInfo entity.User, risk int, imageUrl string) (int, error) {

	err := DefaultOCRService().CheckIdempotent(uInfo.OpenId)
	if err != nil {
		return 0, err
	}
	err = DefaultOCRService().CheckRisk(risk)
	if err != nil {
		return 0, err
	}
	ok, err := NewPointTransactionCountLimitService(context.NewMioContext()).
		CheckLimit(entity.POINT_REDUCE_PLASTIC, uInfo.OpenId)
	if err != nil {
		return 0, err
	}

	valid, result, err := srv.validateReducePlasticImage(imageUrl)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, errno.ErrCommon.WithMessage("不是有效的图片")
	}

	bizId := util.UUID()
	ctx := context.NewMioContext()
	var point int
	//减碳量
	_, err = NewCarbonTransactionService(ctx).Create(api_types.CreateCarbonTransactionDto{
		OpenId: uInfo.OpenId,
		Type:   entity.CARBON_REDUCE_PLASTIC,
		Value:  1,
		Info:   fmt.Sprintf("%s", result),
		BizId:  bizId,
	})
	if err != nil {
		app.Logger.Error("添加骑行更酷减碳量失败", uInfo.OpenId, imageUrl, err)
	}

	//成长体系
	growth_system.GrowthSystemEPPlasticReduction(growthsystemmsg.GrowthSystemParam{
		TaskType:    string(entity.POINT_REDUCE_PLASTIC),
		TaskSubType: string(entity.POINT_REDUCE_PLASTIC),
		UserId:      strconv.FormatInt(uInfo.ID, 10),
		TaskValue:   1,
	})

	if !ok {
		return point, errno.ErrCommon.WithMessage("今日次数以达到上限")
	}
	_, err = NewPointCollectHistoryService(context.NewMioContext()).CreateHistory(CreateHistoryParam{
		OpenId:          uInfo.OpenId,
		TransactionType: entity.POINT_REDUCE_PLASTIC,
		Info:            fmt.Sprintf("reducePlastic=%v", result),
	})
	if err != nil {
		app.Logger.Error("添加环保减塑记录失败", uInfo.OpenId, imageUrl, err)
	}

	point = entity.PointCollectValueMap[entity.POINT_REDUCE_PLASTIC]
	_, err = NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       uInfo.OpenId,
		Type:         entity.POINT_REDUCE_PLASTIC,
		BizId:        bizId,
		ChangePoint:  int64(point),
		AdditionInfo: fmt.Sprintf("{imageUrl=%s}", imageUrl),
	})

	return point, err
}
