package event

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	eevent "mio/internal/pkg/model/entity/event"
	"mio/internal/pkg/service/event"
	"mio/internal/pkg/util/apiutil"
)

var DefaultEventController = EventController{}

type EventController struct{}

func (EventController) GetEventFullDetail(ctx *gin.Context) (gin.H, error) {
	form := GetEventFullDetailForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	info, err := event.DefaultEventService.GetEventFullInfo(form.EventId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"eventInfo": info,
	}, nil
}
func (EventController) GetEventCategoryList(ctx *gin.Context) (gin.H, error) {
	categoryList, err := event.DefaultEventCategoryService.GetEventCategoryList(event.GetEventCategoryListParam{
		OrderBy: entity.OrderByList{eevent.OrderByEventCategorySortDesc},
	})
	if err != nil {
		return nil, err
	}

	infoList := make([]EventCategoryInfo, 0)
	for _, category := range categoryList {
		infoList = append(infoList, EventCategoryInfo{
			EventCategoryId: category.EventCategoryId,
			Title:           category.Title,
			ImageUrl:        category.ImageUrl,
		})
	}
	return gin.H{
		"categoryList": infoList,
	}, nil
}
func (EventController) GetEventList(ctx *gin.Context) (gin.H, error) {
	form := GetEventListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	list, err := event.DefaultEventService.GetEventList(event.GetEventListParam{
		EventCategoryId: form.EventCategoryId,
		OrderBy:         entity.OrderByList{eevent.OrderByEventSortDesc},
	})
	if err != nil {
		return nil, err
	}

	infoList := make([]EventInfo, 0)
	for _, item := range list {
		infoList = append(infoList, EventInfo{
			EventId:           item.EventId,
			EventTemplateType: item.EventTemplateType,
			Title:             item.Title,
			Subtitle:          item.Subtitle,
			CoverImageUrl:     item.CoverImageUrl,
		})
	}

	return gin.H{
		"eventList": infoList,
	}, nil
}
