package service

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
)

const (
	StepToScoreConvertRatio = 60
	ScoreUpperLimit         = 296
)

var DefaultStepService = StepService{repo: repository.DefaultStepRepository}

type StepService struct {
	repo repository.StepRepository
}

func (srv StepService) FindOrCreateStep(userId int64) (*entity.Step, error) {
	userInfo, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if userInfo.ID == 0 {
		return nil, errno.ErrUserNotFound
	}
	step := srv.repo.FindBy(repository.FindStepBy{
		OpenId: userInfo.OpenId,
	})

	if step.ID != 0 {
		return &step, nil
	}

	step = entity.Step{
		OpenId:         userInfo.OpenId,
		Total:          0,
		LastCheckTime:  model.NewTime().StartOfDay(),
		LastCheckCount: 0,
	}
	return &step, srv.repo.Create(&step)
}
