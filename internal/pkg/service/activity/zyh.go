package activity

import (
	"errors"
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
	repoActivity "mio/internal/pkg/repository/activity"
	"mio/internal/pkg/service/srv_types"
)

var DefaultAnswerService = ZyhService{ctx: context.NewMioContext()}

func NewZyhService(ctx *context.MioContext) *ZyhService {
	return &ZyhService{
		ctx:        ctx,
		zyhRepo:    repoActivity.NewZyhRepository(ctx),
		zyhLogRepo: repoActivity.NewZyhLogRepository(ctx),
	}
}

type ZyhService struct {
	ctx        *context.MioContext
	zyhRepo    repoActivity.ZyhRepository
	zyhLogRepo repoActivity.ZyhLogRepository
}

func (srv ZyhService) GetInfoBy(dto srv_types.GetZyhGetInfoByDTO) (entity.Zyh, error) {
	info := srv.zyhRepo.FindBy(repoActivity.FindZyhById{
		Openid: dto.Openid,
		VolId:  dto.VolId,
	})
	return info, nil
}

func (srv ZyhService) Create(dto srv_types.GetZyhGetInfoByDTO) error {
	info, _ := srv.GetInfoBy(dto)
	if info.Id == 0 {
		//入库
		return srv.zyhRepo.Save(&entity.Zyh{
			Openid: dto.Openid,
			VolId:  dto.VolId,
		})
	}
	return errors.New("志愿汇用户已存在")
}
