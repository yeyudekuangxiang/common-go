package repository

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

func BehaviorLogRepository(ctx *context.MioContext) UserBehaviorLogRepository {
	return UserBehaviorLogRepository{ctx: ctx}
}

type UserBehaviorLogRepository struct {
	ctx *context.MioContext
}

func (repo UserBehaviorLogRepository) Save(params *entity.UserBehaviorLog) error {
	return repo.ctx.DB.Save(params).Error
}

func (repo UserBehaviorLogRepository) Updates(params *entity.UserBehaviorLog) error {
	return repo.ctx.DB.Updates(params).Error
}
