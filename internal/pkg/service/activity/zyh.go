package activity

import (
	"mio/internal/pkg/core/context"
	repoActivity "mio/internal/pkg/repository/activity"
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

/*
func (srv ZyhService) GetInfoBy(dto srv_types.DeleteQuestionAnswerDTO) error {
	a := srv.zyhRepo.FindBy(repoActivity.FindZyhBy{
		Openid: "111",
		VolId:  "1111",
	})
	println(a)
    return nil
	//return a
}*/
