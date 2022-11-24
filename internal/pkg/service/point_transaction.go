package service

import (
	"bytes"
	"encoding/json"
	"github.com/jszwec/csvutil"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/oss"
	"mio/internal/pkg/util"
	"time"
)

type PointTransactionService struct {
	ctx  *context.MioContext
	repo *repository.PointTransactionRepository
}

func NewPointTransactionService(ctx *context.MioContext) *PointTransactionService {
	return &PointTransactionService{
		ctx:  ctx,
		repo: repository.NewPointTransactionRepository(ctx),
	}
}

// CreateTransaction 添加发放积分记录
func (srv PointTransactionService) CreateTransaction(param CreatePointTransactionParam) (*entity.PointTransaction, error) {
	transaction := entity.PointTransaction{
		OpenId:         param.OpenId,
		TransactionId:  param.BizId,
		Type:           param.Type,
		Value:          param.Value,
		CreateTime:     model.Time{Time: time.Now()},
		AdditionalInfo: entity.AdditionalInfo(param.AdditionInfo),
		AdminId:        param.AdminId,
		Note:           param.Note,
	}

	return &transaction, srv.repo.Save(&transaction)
}

// GetListBy 查询记录列表
func (srv PointTransactionService) GetListBy(by repository.GetPointTransactionListBy) []entity.PointTransaction {
	return srv.repo.GetListBy(by)
}
func (srv PointTransactionService) FindBy(by repository.FindPointTransactionBy) (*entity.PointTransaction, error) {
	pt := srv.repo.FindBy(by)
	return &pt, nil
}
func (srv PointTransactionService) PagePointRecord(by GetPointTransactionPageListBy) ([]PointRecord, int64, error) {
	recordList := make([]PointRecord, 0)
	isEmptyCondition, openIds, err := srv.filterPointRecordOpenIds(FilterPointRecordOpenIds{
		OpenId:   by.OpenId,
		Nickname: by.Nickname,
		Phone:    by.Phone,
		UserId:   by.UserId,
	})
	if err != nil {
		return recordList, 0, err
	}

	//没有查询到用户
	if !isEmptyCondition && len(openIds) == 0 {
		return recordList, 0, nil
	}

	pointTranList, total := srv.repo.GetPageListBy(repository.GetPointTransactionPageListBy{
		AdminId:         by.AdminId,
		OpenIds:         openIds,
		StartTime:       by.StartTime,
		EndTime:         by.EndTime,
		StartExpireTime: by.StartExpireTime,
		EndExpireTime:   by.EndExpireTime,
		OrderBy:         entity.OrderByList{entity.OrderByPointTranCTDESC},
		Type:            by.Type,
		Types:           by.Types,
		Offset:          by.Offset,
		Limit:           by.Limit,
	})

	pointService := NewPointService(srv.ctx)
	for _, point := range pointTranList {
		user, err := DefaultUserService.GetUserByOpenId(point.OpenId)
		if err != nil {
			return nil, 0, err
		}

		userPoint, err := pointService.FindByOpenId(user.OpenId)
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
			ExpireTime:     timetool.ToTime(point.ExpireTime.Time).Format("2006/01/02 15:04"),
		}
		recordList = append(recordList, pointRecord)
	}
	return recordList, total, nil
}
func (srv PointTransactionService) filterPointRecordOpenIds(by FilterPointRecordOpenIds) (isEmptyCondition bool, openIds []string, err error) {
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
func (srv PointTransactionService) ExportPointTransactionList(adminId int, by ExportPointTransactionListBy) error {
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

		getPageListBy := GetPointTransactionPageListBy{}
		err = util.MapTo(by, &getPageListBy)
		if err != nil {
			return
		}
		getPageListBy.Offset = 0
		getPageListBy.Limit = 100
		type csv struct {
			ID             int64  `csv:"编号"`
			UserId         int64  `csv:"用户ID"`
			OpenId         string `csv:"openId"`
			Phone          string `csv:"用户手机号"`
			Nickname       string `csv:"微信昵称"`
			BalanceOfPoint int64  `csv:"剩余积分"`
			Type           string `csv:"积分变动类型"`
			Value          int64  `csv:"积分变动数量"`
			Time           string `csv:"积分变动时间"`
			Info           string `csv:"附带信息"`
		}
		csvList := make([]csv, 0)
		for {
			list, _, err := srv.PagePointRecord(getPageListBy)
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
			getPageListBy.Offset += getPageListBy.Limit
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
		filePath, err := oss.DefaultOssService.PutObject("images/file-export/point/"+fileName, bytes.NewReader(data))
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
			Url:    oss.DefaultOssService.FullUrl(filePath),
		})
		if err != nil {
			app.Logger.Error("更新导出状态失败", fileExport.ID, err)
		}
		return
	}()
	return nil
}

func (srv PointTransactionService) PageAdjustPointRecord(param GetPointAdjustRecordPageListParam) ([]PointRecord, int64, error) {
	types := make([]entity.PointTransactionType, 0)
	if param.Type != "" {
		types = append(types, param.Type)
	} else {
		types = []entity.PointTransactionType{
			entity.POINT_SYSTEM_REDUCE,
			entity.POINT_SYSTEM_ADD,
		}
	}

	return srv.PagePointRecord(GetPointTransactionPageListBy{
		UserId:    param.UserId,
		AdminId:   param.AdminId,
		Nickname:  param.Nickname,
		OpenId:    param.OpenId,
		Phone:     param.Phone,
		Types:     types,
		StartTime: model.Time{Time: param.StartTime},
		EndTime:   model.Time{Time: param.EndTime},
		Offset:    param.Offset,
		Limit:     param.Limit,
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

func (srv PointTransactionService) CountByToday(openId string, types entity.PointTransactionType) ([]map[string]interface{}, int64, error) {
	return srv.repo.CountByToday(repository.GetPointTransactionCountBy{
		OpenId: openId,
		Type:   types,
	})
}

func (srv PointTransactionService) CountByMonth(openId string, types entity.PointTransactionType) ([]map[string]interface{}, int64, error) {
	return srv.repo.CountByMonth(repository.GetPointTransactionCountBy{
		OpenId: openId,
		Type:   types,
	})
}
