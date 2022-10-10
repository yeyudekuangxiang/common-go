package activity

import (
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
)

func NewZyhLogRepository(ctx *context.MioContext) ZyhLogRepository {
	return ZyhLogRepository{ctx: ctx}
}

type ZyhLogRepository struct {
	ctx *context.MioContext
}

func (repo ZyhLogRepository) Save(zyh *entity.ZyhLog) error {
	return repo.ctx.DB.Save(zyh).Error
}
