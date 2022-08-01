package event

import (
	"mio/internal/pkg/model/entity"
	eevent "mio/internal/pkg/model/entity/event"
	revent "mio/internal/pkg/repository/event"
	"time"
)

var DefaultEventParticipationService = EventParticipationService{repo: revent.DefaultEventParticipationRepository}

type EventParticipationService struct {
	repo revent.EventParticipationRepository
}

// ParticipateEvent  参加活动
func (srv EventParticipationService) ParticipateEvent(user entity.User, list []ParticipateEventParam) error {
	epList := make([]eevent.EventParticipation, 0)
	for _, item := range list {

		ev, err := DefaultEventService.FindEvent(FindEventParam{
			ProductItemId: item.ProductItemId,
		})
		if err != nil {
			return err
		}
		if ev.ID == 0 {
			continue
		}

		//更新参与活动的数量
		_ = DefaultEventService.AddEventParticipationCount(ev.EventId, 1)

		epList = append(epList, eevent.EventParticipation{
			EventId:   ev.EventId,
			Count:     item.Count,
			AvatarUrl: user.AvatarUrl,
			Nickname:  user.Nickname,
			Time:      time.Now(),
		})
	}

	return srv.repo.CreateBatch(&epList)
}
func (srv EventParticipationService) GetParticipationPageList(param GetParticipationPageListParam) ([]eevent.EventParticipation, int64, error) {
	return srv.repo.GetParticipationPageList(revent.GetParticipationPageListBy{
		EventId: param.EventId,
		Limit:   param.Limit,
		Offset:  param.Offset,
		OrderBy: param.OrderBy,
	})
}
