package service

import (
	"bytes"
	"encoding/json"
	"github.com/jszwec/csvutil"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"time"
)

var DefaultPointTransactionService = NewPointTransactionService(repository.DefaultPointTransactionRepository)

func NewPointTransactionService(repo repository.PointTransactionRepository) PointTransactionService {
	return PointTransactionService{
		repo: repo,
	}
}

type PointTransactionService struct {
	repo repository.PointTransactionRepository
}

// Create 添加发放积分记录并且更新用户剩余积分
func (srv PointTransactionService) Create(param CreatePointTransactionParam) (*entity.PointTransaction, error) {

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
		AdminId:        param.AdminId,
		Note:           param.Note,
	}

	err = srv.repo.Save(&transaction)
	if err != nil {
		return nil, err
	}

	_ = DefaultPointService.RefreshBalance(param.OpenId)

	return &transaction, nil
}

// GetListBy 查询记录列表
func (srv PointTransactionService) GetListBy(by repository.GetPointTransactionListBy) []entity.PointTransaction {
	return repository.DefaultPointTransactionRepository.GetListBy(by)
}
func (srv PointTransactionService) FindBy(by repository.FindPointTransactionBy) (*entity.PointTransaction, error) {
	pt := srv.repo.FindBy(by)
	return &pt, nil
}
func (srv PointTransactionService) GetPageListBy(by GetPointTransactionPageListBy) ([]PointRecord, int64, error) {
	recordList := make([]PointRecord, 0)
	isEmptyCondition, openIds, err := srv.getOpenIds(by)
	if err != nil {
		return recordList, 0, err
	}

	//没有查询到用户
	if !isEmptyCondition && len(openIds) == 0 {
		return recordList, 0, nil
	}

	pointTranList, total := srv.repo.GetPageListBy(repository.GetPointTransactionPageListBy{
		OpenIds:   openIds,
		StartTime: by.StartTime,
		EndTime:   by.EndTime,
		OrderBy:   entity.OrderByList{entity.OrderByPointTranCTDESC},
		Type:      by.Type,
		Types:     by.Types,
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

		admin, err := DefaultSystemAdminService.GetAdminById(point.AdminId)
		if err != nil {
			return nil, 0, err
		}

		pointRecord := PointRecord{
			ID:             point.ID,
			User:           *user,
			BalanceOfPoint: userPoint.Balance,
			Type:           point.Type,
			TypeText:       point.Type.RealText(),
			Value:          point.Value,
			CreateTime:     point.CreateTime,
			AdditionalInfo: string(point.AdditionalInfo),
			Note:           point.Note,
			Admin:          *admin,
		}
		recordList = append(recordList, pointRecord)
	}
	return recordList, total, nil
}
func (srv PointTransactionService) getOpenIds(by GetPointTransactionPageListBy) (isEmptyCondition bool, openIds []string, err error) {
	openIds = make([]string, 0)

	if by.UserId == 0 && by.Nickname == "" && by.OpenId == "" && by.Phone == "" {
		isEmptyCondition = true
		return
	}

	userList, err := DefaultUserService.GetUserListBy(repository.GetUserListBy{
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
func (srv PointTransactionService) GetPointTransactionTypeList() []PointTransactionTypeInfo {
	list := make([]PointTransactionTypeInfo, 0)
	for _, t := range entity.PointTransactionTypeList {
		list = append(list, PointTransactionTypeInfo{
			Type:     t,
			TypeText: t.RealText(),
		})
	}
	return list
}
func (srv PointTransactionService) ExportPointTransactionList(adminId int, by GetPointTransactionPageListBy) error {
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
			list, _, err := srv.GetPageListBy(by)
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

		fileName := util.UUID() + ".csv"
		link, err := DefaultOssService.PutObject("static/mp2c/images/file-export/point/"+fileName, bytes.NewReader(data))
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
			Url:    link,
		})
		if err != nil {
			app.Logger.Error("更新导出状态失败", fileExport.ID, err)
		}
		return
	}()
	return nil
}
func (srv PointTransactionService) AdminAdjustUserPoint(adminId int, param AdminAdjustUserPointParam) error {
	user, err := DefaultUserService.GetUserBy(repository.GetUserBy{
		OpenId:     param.OpenId,
		LikeMobile: param.Phone,
	})
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errno.ErrUserNotFound
	}
	_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:  param.OpenId,
		Type:    param.Type,
		Value:   param.Value,
		AdminId: adminId,
		Note:    param.Note,
	})
	return err
}
func (srv PointTransactionService) GetAdjustRecordPageList(param GetPointAdjustRecordPageListParam) ([]PointRecord, int64, error) {
	types := make([]entity.PointTransactionType, 0)
	if param.Type != "" {
		types = append(types, param.Type)
	} else {
		types = []entity.PointTransactionType{
			entity.POINT_SYSTEM_REDUCE,
			entity.POINT_SYSTEM_ADD,
		}
	}
	return DefaultPointTransactionService.GetPageListBy(GetPointTransactionPageListBy{
		OpenId: param.OpenId,
		Phone:  param.Phone,
		Types:  types,
		Offset: param.Offset,
		Limit:  param.Limit,
	})
}

func (srv PointTransactionService) GetAdjustPointTransactionTypeList() []PointTransactionTypeInfo {
	return []PointTransactionTypeInfo{
		{
			Type:     entity.POINT_SYSTEM_REDUCE,
			TypeText: entity.POINT_SYSTEM_REDUCE.RealText(),
		},
		{
			Type:     entity.POINT_SYSTEM_ADD,
			TypeText: entity.POINT_SYSTEM_ADD.RealText(),
		},
	}
}
