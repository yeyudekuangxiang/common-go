package service

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

func BehaviorLogService(ctx *context.MioContext) *UserBehaviorLogService {
	return &UserBehaviorLogService{
		ctx:              ctx,
		BehaviorLogModel: repository.BehaviorLogRepository(ctx),
	}
}

type UserBehaviorLogService struct {
	ctx              *context.MioContext
	BehaviorLogModel repository.UserBehaviorLogRepository
}

func (srv UserBehaviorLogService) Save(params UserBehaviorLogParam) {
	err := srv.BehaviorLogModel.Save(&entity.UserBehaviorLog{
		Tp:   params.Tp,
		Data: params.Data,
		Ip:   params.Ip,
	})
	if err != nil {
		app.Logger.Errorf("入参参数保存失败:%s", err.Error())
	}
	return
}

func (srv UserBehaviorLogService) Updates(params UserBehaviorLogParam) {
	err := srv.BehaviorLogModel.Updates(&entity.UserBehaviorLog{
		Result:     params.Result,
		ResultCode: params.ResultCode,
	})
	if err != nil {
		app.Logger.Errorf("出参参数保存失败:%s", err.Error())
	}
	return
}
