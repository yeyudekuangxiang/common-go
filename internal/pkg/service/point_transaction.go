package service

import (
	"encoding/json"
	"github.com/jszwec/csvutil"
	"io/ioutil"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"time"
)

var DefaultPointTransactionService = NewPointTransactionService(repository2.DefaultPointTransactionRepository)

func NewPointTransactionService(repo repository2.PointTransactionRepository) PointTransactionService {
	return PointTransactionService{
		repo: repo,
	}
}

type PointTransactionService struct {
	repo repository2.PointTransactionRepository
}

// Create 添加发放积分记录并且更新用户剩余积分
func (p PointTransactionService) Create(param CreatePointTransactionParam) (*entity.PointTransaction, error) {

	err := DefaultPointTransactionCountLimitService.CheckLimitAndUpdate(param.Type, param.OpenId)
	if err != nil {
		return nil, err
	}

	transaction := entity.PointTransaction{
		OpenId:         param.OpenId,
		TransactionId:  util.UUID(),
		Type:           param.Type,
		Value:          param.Value,
		CreateTime:     model.Time{Time: time.Now()},
		AdditionalInfo: entity.AdditionalInfo(param.AdditionInfo),
	}

	err = p.repo.Save(&transaction)
	if err != nil {
		return nil, err
	}

	_ = DefaultPointService.RefreshBalance(param.OpenId)

	return &transaction, nil
}

// GetListBy 查询记录列表
func (p PointTransactionService) GetListBy(by repository2.GetPointTransactionListBy) []entity.PointTransaction {
	return repository2.DefaultPointTransactionRepository.GetListBy(by)
}
func (p PointTransactionService) FindBy(by repository2.FindPointTransactionBy) (*entity.PointTransaction, error) {
	pt := p.repo.FindBy(by)
	return &pt, nil
}
func (p PointTransactionService) GetPageListBy(by GetPointTransactionPageListBy) ([]PointRecord, int64, error) {
	recordList := make([]PointRecord, 0)
	isEmptyCondition, openIds, err := p.getOpenIds(by)
	if err != nil {
		return recordList, 0, err
	}

	//没有查询到用户
	if !isEmptyCondition && len(openIds) == 0 {
		return recordList, 0, nil
	}

	pointTranList, total := p.repo.GetPageListBy(repository2.GetPointTransactionPageListBy{
		OpenIds:   openIds,
		StartTime: by.StartTime,
		EndTime:   by.EndTime,
		OrderBy:   entity.OrderByList{entity.OrderByPointTranCTDESC},
		Type:      by.Type,
		Offset:    by.Offset,
		Limit:     by.Limit,
	})

	for _, point := range pointTranList {
		user, err := DefaultUserService.GetUserByOpenId(point.OpenId)
		if err != nil {
			return nil, 0, err
		}

		userPoint, err := DefaultPointService.FindByOpenId(user.OpenId)
		if err != nil {
			return nil, 0, err
		}

		pointRecord := PointRecord{
			ID:             point.ID,
			User:           *user,
			BalanceOfPoint: userPoint.Balance,
			Type:           point.Type,
			TypeText:       point.Type.Text(),
			Value:          point.Value,
			CreateTime:     point.CreateTime,
			AdditionalInfo: string(point.AdditionalInfo),
		}
		recordList = append(recordList, pointRecord)
	}
	return recordList, total, nil
}
func (p PointTransactionService) getOpenIds(by GetPointTransactionPageListBy) (isEmptyCondition bool, openIds []string, err error) {
	openIds = make([]string, 0)

	if by.UserId == 0 && by.Nickname == "" && by.OpenId == "" && by.Phone == "" {
		isEmptyCondition = true
		return
	}

	userList, err := DefaultUserService.GetUserListBy(repository2.GetUserListBy{
		LikeMobile: by.Phone,
		OpenId:     by.OpenId,
		Nickname:   by.Nickname,
		UserId:     by.UserId,
	})
	if err != nil {
		return
	}

	for _, user := range userList {
		openIds = append(openIds, user.OpenId)
	}
	return
}
func (p PointTransactionService) GetPointTransactionTypeList() []PointTransactionTypeInfo {
	list := make([]PointTransactionTypeInfo, 0)
	for _, t := range entity.PointTransactionTypeList {
		list = append(list, PointTransactionTypeInfo{
			Type:     t,
			TypeText: t.Text(),
		})
	}
	return list
}
func (p PointTransactionService) ExportPointTransactionList(adminId int64, by GetPointTransactionPageListBy) error {
	param, err := json.Marshal(by)
	if err != nil {
		return err
	}

	fileExport, err := DefaultFileExportService.Add(AddFileExportParam{
		AdminId: adminId,
		Params:  string(param),
		Type:    entity.FileExportTypePoint,
	})
	if err != nil {
		return err
	}
	go func() {
		_, err := DefaultFileExportService.Update(fileExport.ID, UpdateFileExportParam{
			Status: entity.FileExportStatusProgress,
		})
		if err != nil {
			app.Logger.Error("更新导出状态失败", fileExport.ID, err)
		}
		by.Offset = 0
		by.Limit = 100
		type csv struct {
			ID             int64  `csv:"编号"`
			UserId         int64  `csv:"用户ID"`
			OpenId         string `csv:"openId"`
			Phone          string `csv:"用户手机号"`
			Nickname       string `csv:"微信昵称"`
			BalanceOfPoint int    `csv:"剩余积分"`
			Type           string `csv:"积分变动类型"`
			Value          int    `csv:"积分变动数量"`
			Time           string `csv:"积分变动时间"`
			Info           string `csv:"附带信息"`
		}
		csvList := make([]csv, 0)
		for {
			list, _, err := p.GetPageListBy(by)
			if err != nil {
				_, err := DefaultFileExportService.Update(fileExport.ID, UpdateFileExportParam{
					Status:  entity.FileExportStatusFailed,
					Message: err.Error(),
				})
				if err != nil {
					app.Logger.Error("更新导出状态失败", fileExport.ID, err)
				}
				return
			}
			if len(list) == 0 {
				break
			}
			for _, item := range list {
				csvList = append(csvList, csv{
					ID:             item.ID,
					UserId:         item.User.ID,
					OpenId:         item.User.OpenId,
					Phone:          item.User.PhoneNumber,
					Nickname:       item.User.Nickname,
					BalanceOfPoint: item.BalanceOfPoint,
					Type:           item.TypeText,
					Value:          item.Value,
					Time:           item.CreateTime.String(),
					Info:           item.AdditionalInfo,
				})
			}
			by.Offset += by.Limit
		}

		data, err := csvutil.Marshal(csvList)
		if err != nil {
			_, err := DefaultFileExportService.Update(fileExport.ID, UpdateFileExportParam{
				Status:  entity.FileExportStatusFailed,
				Message: err.Error(),
			})
			if err != nil {
				app.Logger.Error("更新导出状态失败", fileExport.ID, err)
			}
			return
		}

		url := util.UUID() + ".csv"
		err = ioutil.WriteFile(url, data, 0755)
		if err != nil {
			_, err := DefaultFileExportService.Update(fileExport.ID, UpdateFileExportParam{
				Status:  entity.FileExportStatusFailed,
				Message: err.Error(),
			})
			if err != nil {
				app.Logger.Error("更新导出状态失败", fileExport.ID, err)
			}
			return
		}
		_, err = DefaultFileExportService.Update(fileExport.ID, UpdateFileExportParam{
			Status: entity.FileExportStatusSuccess,
			Url:    url,
		})
		if err != nil {
			app.Logger.Error("更新导出状态失败", fileExport.ID, err)
		}
		return
	}()
	return nil
}
