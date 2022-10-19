package activity

import (
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
	"mio/internal/pkg/repository/repotypes"
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

func (repo ZyhLogRepository) GetListBy(by repotypes.GetZyhListBy) ([]entity.ZyhLog, error) {
	list := make([]entity.ZyhLog, 0)
	db := repo.ctx.DB.Model(entity.ZyhLog{})
	if by.Openid != "" {
		db.Where("openid", by.Openid)
	}
	db.Order("id desc")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}
