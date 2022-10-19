package activity

import (
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
)

func NewYtxLogRepository(ctx *context.MioContext) YtxLogRepository {
	return YtxLogRepository{ctx: ctx}
}

type YtxLogRepository struct {
	ctx *context.MioContext
}

func (repo YtxLogRepository) Save(ytx *entity.YtxLog) error {
	return repo.ctx.DB.Save(ytx).Error
}
