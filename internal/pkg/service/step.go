package service

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util/timeutils"
	"mio/pkg/errno"
	"time"
)

const (
	// StepToScoreConvertRatio 步行数量和积分比例  60步等于1积分
	StepToScoreConvertRatio = 60
	// StepScoreUpperLimit 每天步行最多获取的积分数量
	StepScoreUpperLimit = 296
)

// DefaultStepService 默认步行服务
var DefaultStepService = StepService{repo: repository.DefaultStepRepository}

// StepService 步行服务
type StepService struct {
	repo repository.StepRepository
}

// FindOrCreateStep 查询或者创建步行记录
func (srv StepService) FindOrCreateStep(userId int64) (*entity.Step, error) {
	step := srv.repo.FindBy(repository.FindStepBy{
		UserId: userId,
	})

	if step.ID != 0 {
		return &step, nil
	}

	user, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errno.ErrUserNotFound
	}
	step = entity.Step{
		UserId:         userId,
		OpenId:         user.OpenId,
		Total:          0,
		LastCheckTime:  model.NewTime().StartOfDay(),
		LastCheckCount: 0,
	}
	return &step, srv.repo.Create(&step)
}

// UpdateStepTotal 更新用户步行总数
func (srv StepService) UpdateStepTotal(userId int64) error {
	step, err := srv.FindOrCreateStep(userId)
	if err != nil {
		return err
	}

	//查询未统计在内的步行数据历史
	stepHistoryList, err := DefaultStepHistoryService.GetStepHistoryList(GetStepHistoryListBy{
		StartRecordedTime: step.LastSumHistoryTime,
		UserId:            userId,
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
func (srv StepService) WeeklyHistory(userId int64) (*WeeklyHistoryInfo, error) {

	//最近7天记录
	historyList, err := DefaultStepHistoryService.GetStepHistoryList(GetStepHistoryListBy{
		UserId:            userId,
		StartRecordedTime: model.Time{Time: timeutils.StartOfDay(time.Now().Add(-6 * time.Hour * 24))},
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
	step, err := srv.FindOrCreateStep(userId)
	if err != nil {
		return nil, err
	}

	//历史总步数和总天数
	total, days := DefaultStepHistoryService.GetUserLifeStepInfo(userId)

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
