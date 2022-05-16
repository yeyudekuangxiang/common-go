package service

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
)

var DefaultEventParticipationService = EventParticipationService{}

type EventParticipationService struct {
	repo repository.EventParticipationRepository
}

// ParticipateEvent  参加活动
func (srv EventParticipationService) ParticipateEvent(userId int64, list []ParticipateEventParam) error {
	user, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errno.ErrUserNotFound
	}

	epList := make([]entity.EventParticipation, 0)
	for _, item := range list {

		event, err := DefaultEventService.FindEvent(FindEventBy{
			ProductItemId: item.ProductItemId,
		})
		if err != nil {
			return err
		}
		if event.ID == 0 {
			continue
		}

		epList = append(epList, entity.EventParticipation{
			EventId:   event.EventId,
			Count:     item.Count,
			AvatarUrl: user.AvatarUrl,
			Nickname:  user.Nickname,
			Time:      model.NewTime(),
		})
	}
	return srv.repo.CreateBatch(&epList)
}
