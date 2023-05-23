package event

import (
	eevent "mio/internal/pkg/model/entity/event"
	revent "mio/internal/pkg/repository/event"
)

var DefaultEventCategoryService = EventCategoryService{repo: revent.DefaultEventCategoryRepository}

type EventCategoryService struct {
	repo revent.EventCategoryRepository
}

func (srv EventCategoryService) GetEventCategoryList(param GetEventCategoryListParam) ([]eevent.APIEventCategory, error) {
	return srv.repo.GetEventCategoryList(revent.GetEventCategoryListBy{
		OrderBy: param.OrderBy,
		Active:  param.Active,
		Display: param.Display,
	})
}
