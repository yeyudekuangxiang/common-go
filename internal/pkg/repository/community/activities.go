package community

import (
	mioContext "mio/internal/pkg/core/context"
)

type (
	ActivitiesModel interface {
	}

	defaultCommunityActivitiesModel struct {
		ctx *mioContext.MioContext
	}
)

func NewCommunityActivitiesModel(ctx *mioContext.MioContext) ActivitiesModel {
	return defaultCommunityActivitiesModel{
		ctx: ctx,
	}
}
