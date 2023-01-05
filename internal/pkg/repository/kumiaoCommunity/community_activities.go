package kumiaoCommunity

import (
	mioContext "mio/internal/pkg/core/context"
)

type (
	CommunityActivitiesModel interface {
	}

	defaultCommunityActivitiesModel struct {
		ctx *mioContext.MioContext
	}
)

func NewCommunityActivitiesModel(ctx *mioContext.MioContext) CommunityActivitiesModel {
	return defaultCommunityActivitiesModel{
		ctx: ctx,
	}
}
