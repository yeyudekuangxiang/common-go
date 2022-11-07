package service

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

type (
	BdScenePrePointService interface {
		Save(data *entity.BdScenePrePoint) error
	}

	defaultBdScenePrePointService struct {
		ctx           *mioContext.MioContext
		sceneUser     repository.BdSceneUserModel
		prePointModel repository.BdScenePrePointModel
	}
)

func (srv defaultBdScenePrePointService) Save(data *entity.BdScenePrePoint) error {
	return srv.prePointModel.Save(data)
}

func NewBdScenePrePointService(ctx *mioContext.MioContext) BdScenePrePointService {
	return &defaultBdScenePrePointService{
		ctx:           ctx,
		sceneUser:     repository.NewBdSceneUserModel(ctx),
		prePointModel: repository.NewBdScenePrePointModel(ctx),
	}
}
