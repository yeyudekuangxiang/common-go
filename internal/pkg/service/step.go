package service

import (
	contextRedis "context"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/commontool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"mio/config"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strconv"
	"time"
)

const (
	// StepToScoreConvertRatio 步行数量和积分比例  90步等于1积分
	StepToScoreConvertRatio = 90
	// StepScoreUpperLimit 每天步行最多获取的积分数量
	StepScoreUpperLimit = 133
)

// DefaultStepService 默认步行服务
var DefaultStepService = StepService{repo: repository.DefaultStepRepository}

// StepService 步行服务
type StepService struct {
	repo repository.StepRepository
}

// FindOrCreateStep 查询或者创建步行记录
func (srv StepService) FindOrCreateStep(openId string) (*entity.Step, error) {
	step := srv.repo.FindBy(repository.FindStepBy{
		OpenId: openId,
	})

	if step.ID != 0 {
		return &step, nil
	}

	step = entity.Step{
		OpenId:         openId,
		Total:          0,
		LastCheckTime:  model.NewTime().StartOfDay(),
		LastCheckCount: 0,
	}
	return &step, srv.repo.Create(&step)
}

// UpdateStepTotal 更新用户步行总数
func (srv StepService) UpdateStepTotal(openId string) error {
	step, err := srv.FindOrCreateStep(openId)
	if err != nil {
		return err
	}

	//查询未统计在内的步行数据历史
	stepHistoryList, err := DefaultStepHistoryService.GetStepHistoryList(GetStepHistoryListBy{
		StartRecordedTime: model.Time{Time: step.LastSumHistoryTime.Add(time.Nanosecond)},
		OpenId:            openId,
		OrderBy:           entity.OrderByList{entity.OrderByStepHistoryTimeDesc},
	})

	if err != nil {
		return err
	}

	total := 0
	for i, item := range stepHistoryList {
		if i == 0 {
			step.LastSumHistoryNum = item.Count
			step.LastSumHistoryTime = item.RecordedTime
		}
		if item.RecordedTime.Equal(step.LastSumHistoryTime.Time) {
			total += item.Count - step.LastSumHistoryNum
		} else {
			total += item.Count
		}
	}

	step.Total += int64(total)
	return srv.repo.Save(step)
}

// WeeklyHistory 获取最近一周
func (srv StepService) WeeklyHistory(openId string) (*WeeklyHistoryInfo, error) {

	//最近7天记录
	historyList, err := DefaultStepHistoryService.GetStepHistoryList(GetStepHistoryListBy{
		OpenId:            openId,
		StartRecordedTime: model.Time{Time: timetool.StartOfDay(time.Now().Add(-6 * time.Hour * 24))},
	})
	if err != nil {
		return nil, err
	}

	//最近7天减少co2
	totalStep := 0
	for _, history := range historyList {
		totalStep += history.Count
	}
	sevenDayCo2 := DefaultCarbonNeutralityService.calculateCO2ByStep(int64(totalStep))

	//查询历史总步数
	step, err := srv.FindOrCreateStep(openId)
	if err != nil {
		return nil, err
	}

	//历史总步数和总天数
	total, days := DefaultStepHistoryService.GetUserLifeStepInfo(openId)

	weeklyHistoryInfo := WeeklyHistoryInfo{}
	weeklyHistoryInfo.SevenDaysCo2 = sevenDayCo2
	weeklyHistoryInfo.LifeSteps = step.Total
	weeklyHistoryInfo.LifeSavedCo2 = DefaultCarbonNeutralityService.calculateCO2ByStep(step.Total)

	if days != 0 {
		weeklyHistoryInfo.AveragePerWeeklyCo2 = DefaultCarbonNeutralityService.calculateCO2ByStep((total * 7) / days)
	}

	weeklyHistoryInfo.StepList = srv.formatWeeklyHistoryStepList(historyList)
	return &weeklyHistoryInfo, nil
}
func (srv StepService) formatWeeklyHistoryStepList(list []entity.StepHistory) []WeeklyHistoryStep {
	stepList := make([]WeeklyHistoryStep, 0)
	for _, history := range list {
		stepList = append(stepList, WeeklyHistoryStep{
			Count:     history.Count,
			Time:      history.RecordedTime,
			Timestamp: history.RecordedEpoch,
		})
	}
	return stepList
}

// RedeemPointFromPendingSteps 领取步行积分
func (srv StepService) RedeemPointFromPendingSteps(openId string, bizId string) (int, error) {
	if !util.DefaultLock.Lock(fmt.Sprintf("RedeemPointFromPendingSteps%s", openId), time.Second*5) {
		return 0, errno.ErrCommon.WithMessage("操作频繁,请稍后再试")
	}
	defer util.DefaultLock.UnLock(fmt.Sprintf("RedeemPointFromPendingSteps%s", openId))

	stepHistory, err := DefaultStepHistoryService.FindStepHistory(FindStepHistoryBy{
		OpenId: openId,
		Day:    model.NewTime().StartOfDay(),
	})

	if err != nil {
		return 0, err
	}

	pendingPoint, pendingStep, err := srv.ComputePendingPoint(openId)
	if err != nil {
		return 0, err
	}

	if pendingPoint == 0 {
		return 0, nil
	}

	step, err := srv.FindOrCreateStep(openId)
	if err != nil {
		return 0, err
	}

	step.LastCheckCount = srv.computeLastCheckedSteps(step.LastCheckTime.Time, step.LastCheckCount) + pendingStep
	step.LastCheckTime = stepHistory.RecordedTime
	err = srv.repo.Save(step)
	if err != nil {
		return 0, err
	}

	_, err = NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		ChangePoint:  pendingPoint,
		BizId:        bizId,
		Type:         entity.POINT_STEP,
		OpenId:       openId,
		AdditionInfo: fmt.Sprintf("{time=%v, count=%d, point=%d}", time.Now(), stepHistory.Count, pendingPoint),
	})
	return int(pendingPoint), err
}

// RedeemCarbonFromPendingSteps 领取步行碳量
func (srv StepService) RedeemCarbonFromPendingSteps(openId string, uid int64, ip string, bizId string) (float64, error) {
	if !util.DefaultLock.Lock(fmt.Sprintf("RedeemCarbonFromPendingSteps%s", openId), time.Second*5) {
		return 0, errno.ErrCommon.WithMessage("操作频繁,请稍后再试")
	}
	defer util.DefaultLock.UnLock(fmt.Sprintf("RedeemCarbonFromPendingSteps%s", openId))

	stepHistory, err := DefaultStepHistoryService.FindStepHistory(FindStepHistoryBy{
		OpenId: openId,
		Day:    model.NewTime().StartOfDay(),
	})
	allStep := decimal.NewFromFloat(float64(stepHistory.Count)) //当天总步数
	UserIdString := strconv.FormatInt(uid, 10)
	todayDate := time.Now().Format("20060102") //当天时间 年月日
	redisKey := fmt.Sprintf(config.RedisKey.UserCarbonStep, todayDate)

	nowTep := app.Redis.ZScore(contextRedis.Background(), redisKey, UserIdString) //已被转化碳的步数
	addStep := allStep.Sub(decimal.NewFromFloat(nowTep.Val()))                    //本次增加的
	addStepFloat, _ := addStep.Float64()
	errRedis := app.Redis.ZIncrBy(contextRedis.Background(), redisKey, addStepFloat, UserIdString).Err()
	if errRedis != nil {
		return 0, errno.ErrCommon.WithMessage("步数更新失败")
	}
	//发碳
	//int64转float64
	carbon, errCarbon := NewCarbonTransactionService(context.NewMioContext()).Create(api_types.CreateCarbonTransactionDto{
		OpenId:  openId,
		UserId:  uid,
		Type:    entity.CARBON_STEP,
		Value:   addStepFloat,
		Info:    fmt.Sprintf("{time=%v, count=%d,allCount=%f}", time.Now(), stepHistory.Count, addStepFloat),
		AdminId: 0,
		Ip:      ip,
		BizId:   bizId,
	})
	if errCarbon != nil {
		return 0, errCarbon
	}
	return carbon, err
}

func (srv StepService) computeLastCheckedSteps(lastCheckedTime time.Time, lastCheckedCount int) int {
	//如果最后一次领积分时间为0 或者 最后一次领取时间不等于今天的开始时间
	if lastCheckedTime.IsZero() || !lastCheckedTime.Equal(model.NewTime().StartOfDay().Time) {
		return 0
	}
	return lastCheckedCount
}

// computePendingStep 计算未领取积分的步行步数
func (srv StepService) computePendingStep(history entity.StepHistory, step entity.Step) int {
	stepUpperLimit := StepToScoreConvertRatio * StepScoreUpperLimit
	lastCheckedSteps := srv.computeLastCheckedSteps(step.LastCheckTime.Time, step.LastCheckCount)

	//之前的领取已经超出了
	if lastCheckedSteps >= stepUpperLimit {
		return 0
	}

	currentStep := commontool.Ternary(history.Count < stepUpperLimit, history.Count, stepUpperLimit).Int()

	result := currentStep - lastCheckedSteps

	return commontool.Ternary(result > 0, result, 0).Int()
}

// computePendingStep 计算可领取的积分数量 和 步行数量
func (srv StepService) computePendingPoint(history entity.StepHistory, step entity.Step) (point int64, stepNum int) {
	pendingStep := srv.computePendingStep(history, step)
	pendingPoint := pendingStep / StepToScoreConvertRatio
	return int64(pendingPoint), pendingPoint * StepToScoreConvertRatio
}

func (srv StepService) ComputePendingPoint(openId string) (point int64, step int, err error) {

	stepHistory, err := DefaultStepHistoryService.FindStepHistory(FindStepHistoryBy{
		OpenId: openId,
		Day:    model.NewTime().StartOfDay(),
	})
	if err != nil {
		return 0, 0, err
	}
	userStep, err := DefaultStepService.FindOrCreateStep(openId)
	if err != nil {
		return 0, 0, err
	}

	pointNum, stepNum := srv.computePendingPoint(*stepHistory, *userStep)
	return pointNum, stepNum, nil
}
