package activity

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
)

func NewZyhRepository(ctx *context.MioContext) ZyhRepository {
	return ZyhRepository{ctx: ctx}
}

type ZyhRepository struct {
	ctx *context.MioContext
}

func (repo ZyhRepository) Save(carbon *entity.Zyh) error {
	return repo.ctx.DB.Save(carbon).Error
}

func (repo ZyhRepository) FindBy(by FindZyhById) entity.Zyh {
	zyh := entity.Zyh{}
	db := repo.ctx.DB.Model(zyh)
	if by.Openid != "" {
		db.Where("openid = ?", by.Openid)
	}
	if by.VolId != "" {
		db.Where("vol_id = ?", by.VolId)
	}
	if err := db.First(&zyh).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return zyh
}
