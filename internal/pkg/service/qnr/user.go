package question

import (
	"mio/internal/pkg/core/context"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	repoQnr "mio/internal/pkg/repository/qnr"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
)

var DefaultUserService = UserService{ctx: context.NewMioContext()}

func NewUserService(ctx *context.MioContext) *UserService {
	return &UserService{
		ctx:  ctx,
		repo: repoQnr.NewUserRepository(ctx),
	}
}

type UserService struct {
	ctx  *context.MioContext
	repo *repoQnr.UserRepository
}

func (srv UserService) GetUserInfo(dto srv_types.GetQnrUserDTO) (qnrEntity.User, error) {
	list := srv.repo.FindBy(repotypes.GetQuestUserGetById{UserId: dto.UserId})
	return list, nil
}
