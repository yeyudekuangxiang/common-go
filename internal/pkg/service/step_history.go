package service

import (
	"github.com/medivhzhan/weapp/v3"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
	"time"
)

// DefaultStepHistoryService 默认步行历史服务
var DefaultStepHistoryService = StepHistoryService{repo: repository.DefaultStepHistoryRepository}

// StepHistoryService 步行历史服务
type StepHistoryService struct {
	repo repository.StepHistoryRepository
}

// FindStepHistory 查询一条步行历史记录
func (srv StepHistoryService) FindStepHistory(by FindStepHistoryBy) (*entity.StepHistory, error) {
	step := srv.repo.FindBy(repository.FindStepHistoryBy{
		OpenId:  by.OpenId,
		Day:     by.Day,
		OrderBy: by.OrderBy,
	})
	return &step, nil
}

// UpdateStepHistoryByEncrypted 根据微信运动加密数据更新用户步行历史记录
func (srv StepHistoryService) UpdateStepHistoryByEncrypted(param UpdateStepHistoryByEncryptedParam) error {
	sessionKey, err := DefaultSessionService.MustGetSessionKey(param.OpenId)
	if err != nil {
		app.Logger.Errorf("查询sessionKey失败 %v %+v", param.OpenId, err)
		return err
	}

	runData, err := app.Weapp.DecryptRunData(sessionKey, param.EncryptedData, param.IV)
	if err != nil {
		app.Logger.Errorf("解析微信运动数据失败 %v %+v %+v %+v", param.OpenId, sessionKey, param, err)
		return errno.ErrAuth
	}

	err = srv.updateStepHistoryByList(param.OpenId, runData.StepInfoList)
	if err != nil {
		return err
	}

	return DefaultStepService.UpdateStepTotal(param.OpenId)
}

func (srv StepHistoryService) UpdateStepHistoryByList(openId string, stepInfoList []weapp.SetpInfo) error {
	return srv.updateStepHistoryByList(openId, stepInfoList)
}

// updateStepHistoryByList 根据微信运动数据列表创建或者更新步行历史记录(最多更新最近8天数据)
func (srv StepHistoryService) updateStepHistoryByList(openId string, stepInfoList []weapp.SetpInfo) error {

	//只更新最近8天的
	if len(stepInfoList) > 8 {
		stepInfoList = stepInfoList[len(stepInfoList)-8:]
	}

	now := timetool.Now().StartOfDay().Unix()
	for _, stepInfo := range stepInfoList {
		//判断是否是今天
		if stepInfo.Timestamp == now {
			err := srv.createOrUpdateTodayStep(openId, stepInfo)
			if err != nil {
				return err
			}
		} else if stepInfo.Timestamp < now {
			err := srv.createHistoryStep(openId, stepInfo)
			if err != nil {
				app.Logger.Errorf("更新历史步数失败 %s %+v", openId, stepInfo)
			}
		}
	}
	return nil
}
func (srv StepHistoryService) createOrUpdateTodayStep(openId string, info weapp.SetpInfo) error {
	history := srv.repo.FindBy(repository.FindStepHistoryBy{
		OpenId:        openId,
		RecordedEpoch: info.Timestamp,
	})
	if history.ID == 0 {
		history = entity.StepHistory{
			OpenId:        openId,
			Count:         info.Step,
			RecordedTime:  model.Time{Time: time.Unix(info.Timestamp, 0)},
			RecordedEpoch: info.Timestamp,
		}
		return srv.repo.Create(&history)
	}
	if info.Step > history.Count {
		history.Count = info.Step
	}
	return srv.repo.Save(&history)
}
func (srv StepHistoryService) createHistoryStep(openId string, info weapp.SetpInfo) error {
	history := srv.repo.FindBy(repository.FindStepHistoryBy{
		OpenId:        openId,
		RecordedEpoch: info.Timestamp,
	})
	if history.ID == 0 {
		history = entity.StepHistory{
			OpenId:        openId,
			Count:         info.Step,
			RecordedTime:  model.Time{Time: time.Unix(info.Timestamp, 0)},
			RecordedEpoch: info.Timestamp,
		}
		return srv.repo.Create(&history)
	}
	return nil
}

// GetStepHistoryList 获取步行历史记录列表
func (srv StepHistoryService) GetStepHistoryList(by GetStepHistoryListBy) ([]entity.StepHistory, error) {
	list := srv.repo.GetStepHistoryList(repository.GetStepHistoryListBy{
		OpenId:            by.OpenId,
		RecordedEpochs:    by.RecordEpochs,
		StartRecordedTime: by.StartRecordedTime,
		EndRecordedTime:   by.EndRecordedTime,
		OrderBy:           entity.OrderByList{entity.OrderByStepHistoryTimeDesc},
	})
	return list, nil
}

// GetStepHistoryPageList 获取步行历史记录分页列表
func (srv StepHistoryService) GetStepHistoryPageList(by GetStepHistoryPageListBy) ([]entity.StepHistory, int64, error) {
	list, total := srv.repo.GetStepHistoryPageList(repository.GetStepHistoryPageListBy{
		GetStepHistoryListBy: repository.GetStepHistoryListBy{
			OpenId:            by.OpenId,
			RecordedEpochs:    by.RecordEpochs,
			StartRecordedTime: by.StartRecordedTime,
			EndRecordedTime:   by.EndRecordedTime,
			OrderBy:           entity.OrderByList{entity.OrderByStepHistoryTimeDesc},
		},
		Offset: by.Offset,
		Limit:  by.Limit,
	})
	return list, total, nil
}

// GetUserLifeStepInfo 获取用户历史总步数及总天数
func (srv StepHistoryService) GetUserLifeStepInfo(openId string) (steps int64, days int64) {
	return srv.repo.GetUserLifeStepInfo(openId)
}
