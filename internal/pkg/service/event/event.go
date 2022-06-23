package event

import (
	"database/sql"
	"errors"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/event"
	revent "mio/internal/pkg/repository/event"
	"mio/internal/pkg/service/product"
	"strconv"
	"strings"
	"time"
)

var DefaultEventService = EventService{repo: revent.DefaultEventRepository}

type EventService struct {
	repo revent.EventRepository
}

func (srv EventService) FindEvent(param FindEventParam) (*event.Event, error) {
	ev, err := srv.repo.FindEvent(revent.FindEventBy{
		ProductItemId: param.ProductItemId,
		EventId:       param.EventId,
		Active:        param.Active,
	})
	if err != nil {
		return nil, err
	}
	return &ev, nil
}
func (srv EventService) GetEventFullInfo(eventId string) (*EventFullInfo, error) {
	ev, err := srv.FindEvent(FindEventParam{
		EventId: eventId,
		Active:  sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	if ev.ID == 0 {
		return nil, errors.New("项目不存在")
	}

	participationList, _, err := DefaultEventParticipationService.GetParticipationPageList(GetParticipationPageListParam{
		EventId: eventId,
		Limit:   3,
	})
	if err != nil {
		return nil, err
	}
	participationInfoList := make([]ParticipationInfo, 0)
	for _, participation := range participationList {
		timeStr := ""
		day := int(time.Now().Sub(participation.Time).Hours() / 24)
		if day == 0 {
			timeStr = "今天"
		} else {
			timeStr = fmt.Sprintf("%d天前", day)
		}
		participationInfoList = append(participationInfoList, ParticipationInfo{
			Nickname: participation.Nickname,
			Avatar:   participation.AvatarUrl,
			TimeStr:  timeStr,
			Message:  strings.ReplaceAll(ev.ParticipationTitle, "{X}", strconv.Itoa(participation.Count)),
		})
	}

	setting, err := DefaultEventTemplateService.ParseSetting(ev.EventTemplateType, ev.TemplateSetting)
	if err != nil {
		app.Logger.Error(ev.EventTemplateType, ev.TemplateSetting, err)
		return nil, errors.New("系统异常,请稍后再试")
	}

	eventDetail, err := DefaultEventDetailService.GetFormattedEventDetail(eventId)
	if err != nil {
		return nil, err
	}
	eventRule, err := DefaultEventRuleService.GetFormattedEventRule(eventId)
	if err != nil {
		return nil, err
	}

	productItem, err := product.DefaultProductItemService.FindProductByItemId(ev.ProductItemId)
	if err != nil {
		return nil, err
	}

	return &EventFullInfo{
		EventId:               ev.EventId,
		EventTemplateType:     ev.EventTemplateType,
		Title:                 ev.Title,
		SubTitle:              ev.Subtitle,
		CoverImageUrl:         ev.CoverImageUrl,
		StartTime:             ev.StartTime,
		EndTime:               ev.EndTime,
		ParticipationCount:    ev.ParticipationCount,
		ParticipationSubtitle: ev.ParticipationSubtitle,
		Tags:                  ev.Tag,
		TemplateSetting:       map[string]event.EventTemplateSettingInfo{string(ev.EventTemplateType): setting},
		ParticipationList:     participationInfoList,
		EventDetail:           eventDetail,
		EventRule:             eventRule,
		Cost:                  productItem.Cost,
		Stock:                 productItem.RemainingCount,
	}, nil
}
func (srv EventService) AddEventParticipationCount(eventId string, count int) error {
	ev, err := srv.repo.FindEvent(revent.FindEventBy{
		EventId: eventId,
	})
	if err != nil {
		return err
	}
	if ev.ID != 0 {
		return errors.New("未查询到项目信息")
	}
	ev.ParticipationCount += count
	return srv.repo.Save(&ev)
}
func (srv EventService) GetEventList(param GetEventListParam) ([]event.Event, error) {
	return srv.repo.GetEventList(revent.GetEventListBy{
		EventCategoryId: param.EventCategoryId,
		OrderBy:         param.OrderBy,
		Active:          param.Active,
	})
}
func (srv EventService) GetEventShortInfoList(param GetEventListParam) ([]EventShortInfo, error) {
	eventList, err := srv.repo.GetEventList(revent.GetEventListBy{
		EventCategoryId: param.EventCategoryId,
		OrderBy:         param.OrderBy,
		Active:          sql.NullBool{Bool: true, Valid: true},
	})
	fmt.Printf("%+v\n", eventList)
	if err != nil {
		return nil, err
	}

	productItemIds := make([]string, 0)
	for _, ev := range eventList {
		productItemIds = append(productItemIds, ev.ProductItemId)
	}
	productItems := product.DefaultProductItemService.GetListBy(product.GetProductItemListParam{
		ItemIds: productItemIds,
	})
	productItemMap := product.DefaultProductItemService.ListToMap(productItems)

	eventInfoList := make([]EventShortInfo, 0)

	for _, ev := range eventList {
		eventInfoList = append(eventInfoList, EventShortInfo{
			EventId:           ev.EventId,
			EventTemplateType: ev.EventTemplateType,
			Title:             ev.Title,
			Subtitle:          ev.Subtitle,
			CoverImageUrl:     ev.CoverImageUrl,
			Cost:              productItemMap[ev.ProductItemId].Cost,
		})
	}

	return eventInfoList, nil
}
