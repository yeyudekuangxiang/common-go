package service

import (
	"errors"
	"fmt"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/timeutils"
	"time"
)

var (
	CheckInPointRule = []int{
		100, //第一天
		150, //第二天
		150, //第三天
		200, //第四天
		200, //四五天
		200, //四六天
		300, //第七天
	}
)
var DefaultCheckinService = CheckinService{}

type CheckinService struct {
}

func (srv CheckinService) GetCheckInfo(openId string) (*CheckinInfo, error) {
	history, err := DefaultCheckinHistoryService.FindLastCheckinHistory(openId)
	if err != nil {
		return nil, err
	}
	isCheckedToday := srv.isCheckToday(history.Time.Time)

	checkinInfo := CheckinInfo{
		IsChecked: isCheckedToday,
	}

	checkinInfo.CheckedNumber = history.CheckedNumber

	if isCheckedToday {
		checkinInfo.TodayCheckIndex = history.CheckedNumber
	} else {
		checkinInfo.TodayCheckIndex = srv.nextCheckinDay(history.CheckedNumber)
	}

	checkinInfo.QuickCheckin = srv.getQuickCheckin(checkinInfo.CheckedNumber)
	checkinInfo.Rule = CheckInPointRule
	return &checkinInfo, nil
}
func (srv CheckinService) isCheckToday(lastCheckTime time.Time) bool {
	return timeutils.StartOfDay(lastCheckTime).Equal(timeutils.StartOfDay(time.Now()))
}
func (srv CheckinService) nextCheckinDay(lastCheckNum int) int {
	return lastCheckNum%7 + 1
}
func (srv CheckinService) getQuickCheckin(currentCheckedDay int) []int {
	if currentCheckedDay < 3 {
		return CheckInPointRule[:5]
	} else {
		return CheckInPointRule[2:]
	}
}
func (srv CheckinService) getCheckDayPoint(dayNumber int) int {
	return CheckInPointRule[dayNumber-1]
}
func (srv CheckinService) Checkin(openId string) (int, error) {
	if !util.DefaultLock.Lock(fmt.Sprintf("Checkin%s", openId), time.Second*5) {
		return 0, errors.New("操作太频繁,请稍后再试")
	}

	history, err := DefaultCheckinHistoryService.FindLastCheckinHistory(openId)
	if err != nil {
		return 0, err
	}

	if srv.isCheckToday(history.Time.Time) {
		return history.CheckedNumber, nil
	}

	currentCheckInDay := srv.nextCheckinDay(history.CheckedNumber)

	_, err = DefaultCheckinHistoryService.CreateCheckinHistory(openId, currentCheckInDay)
	if err != nil {
		return 0, err
	}

	point := srv.getCheckDayPoint(currentCheckInDay)

	_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId: openId,
		Type:   entity.POINT_CHECK_IN,
		Value:  point,
	})

	return currentCheckInDay, err
}
