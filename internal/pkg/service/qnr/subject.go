package question

import (
	"mio/internal/pkg/core/context"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	repoQnr "mio/internal/pkg/repository/qnr"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
)

var DefaultSubjectService = SubjectService{ctx: context.NewMioContext()}

func NewSubjectService(ctx *context.MioContext) *SubjectService {
	return &SubjectService{
		ctx:  ctx,
		repo: repoQnr.NewSubjectRepository(ctx),
	}
}

type SubjectService struct {
	ctx  *context.MioContext
	repo *repoQnr.SubjectRepository
}

func (srv SubjectService) GetPageList(dto srv_types.GetQnrSubjectDTO) ([]qnrEntity.Subject, error) {
	list, err := srv.repo.List(repotypes.GetQuestSubjectGetListBy{QnrId: dto.QnrId})
	if err != nil {
		return nil, err
	}
	return list, nil
}
