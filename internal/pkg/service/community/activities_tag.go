package community

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/community"
	"mio/pkg/errno"
)

type (
	CommunityActivitiesTagService interface {
		List() ([]entity.CommunityActivitiesTag, error)
		GetPageList(param community.FindAllActivitiesTagParams) ([]entity.CommunityActivitiesTag, int64, error)
		GetOne(id int64) (entity.CommunityActivitiesTag, error)
	}

	defaultCommunityActivitiesTagService struct {
		ctx      *mioContext.MioContext
		tagModel community.ActivitiesTagModel
	}
)

func NewCommunityActivitiesTagService(ctx *mioContext.MioContext) CommunityActivitiesTagService {
	return defaultCommunityActivitiesTagService{
		ctx:      ctx,
		tagModel: community.NewCommunityActivitiesTagModel(ctx),
	}
}

func (srv defaultCommunityActivitiesTagService) List() ([]entity.CommunityActivitiesTag, error) {
	tags, err := srv.tagModel.List()
	if err != nil {
		return []entity.CommunityActivitiesTag{}, errno.ErrInternalServer.WithMessage(err.Error())
	}
	return tags, nil
}

func (srv defaultCommunityActivitiesTagService) GetPageList(param community.FindAllActivitiesTagParams) ([]entity.CommunityActivitiesTag, int64, error) {
	list, total, err := srv.tagModel.GetPageList(param)
	if err != nil {
		return nil, 0, errno.ErrInternalServer.WithMessage(err.Error())
	}
	return list, total, nil
}

func (srv defaultCommunityActivitiesTagService) GetOne(id int64) (entity.CommunityActivitiesTag, error) {
	tag, err := srv.tagModel.GetById(id)
	if err != nil {
		return entity.CommunityActivitiesTag{}, errno.ErrInternalServer.WithMessage(err.Error())
	}
	if tag.Id == 0 {
		return entity.CommunityActivitiesTag{}, errno.ErrCommon.WithMessage("未找到该标签")
	}
	return tag, nil
}
