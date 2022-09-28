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
		ctx:     ctx,
		zyhRepo: repoActivity.NewZyhRepository(ctx),
	}
}

type ZyhService struct {
	ctx     *context.MioContext
	zyhRepo repoActivity.ZyhRepository
}

func (srv ZyhService) GetInfoBy(dto srv_types.GetZyhGetInfoByDTO) (entity.Zyh, error) {
	info := srv.zyhRepo.FindBy(repoActivity.FindZyhById{
		Openid: dto.Openid,
		VolId:  dto.VolId,
	})
	if info.Id == 0 {
		return entity.Zyh{}, errors.New("无用户信息")
	}
	return info, nil
}
