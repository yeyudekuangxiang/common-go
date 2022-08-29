package question

import (
	"mio/internal/pkg/core/context"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	repoQnr "mio/internal/pkg/repository/qnr"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
)

var DefaultOptionService = OptionService{ctx: context.NewMioContext()}

func NewOptionService(ctx *context.MioContext) *OptionService {
	return &OptionService{
		ctx:  ctx,
		repo: repoQnr.NewOptionRepository(ctx),
	}
}

type OptionService struct {
	ctx  *context.MioContext
	repo *repoQnr.OptionRepository
}

func (srv OptionService) GetPageList(dto srv_types.GetQnrOptionDTO) ([]qnrEntity.Option, error) {
	list, err := srv.repo.GetListBy(repotypes.GetQuestOptionGetListBy{SubjectIds: dto.SubjectIds})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (srv OptionService) CreateInBatches(dto []qnrEntity.Option) error {
	err := srv.repo.CreateInBatches(dto)
	return err
}
