package activity

import (
	"fmt"
	"github.com/pkg/errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	entity2 "mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/activity"
	activity2 "mio/internal/pkg/repository/activity"
	service2 "mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"time"
)

var (
	GMProductItemId  = "146bc6b0-bea8-45ba-b75a-e0df19f3b5ca"
	GMNewUserTime, _ = time.Parse("2006-01-02 15:04:05", "2022-04-25 15:00:00")
	GMEndTime, _     = time.Parse("2006-01-02 15:04:05", "2022-03-25 15:00:00")
)
var DefaultGMService = GMService{}

type GMService struct {
}

func (srv GMService) Order(userId int64, addressId string) (*entity2.Order, error) {
	record, err := srv.FindOrCreateGMRecord(userId)
	if err != nil {
		return nil, err
	}
	if record.PrizeStatus == 1 {
		return nil, errors.New("未完成挑战,请完成挑战后再领取")
	}
	if record.PrizeStatus == 3 {
		return nil, errors.New("已领取过奖励,无法继续领取	")
	}
	if record.PrizeStatus == 2 {
		record.PrizeStatus = 3
		err = activity2.DefaultGMRecordRepository.Save(record)
		if err != nil {
			return nil, err
		}
		order, err := service2.DefaultOrderService.SubmitOrderForGreenMonday(service2.SubmitOrderForGreenParam{
			AddressId: addressId,
			UserId:    userId,
			ItemId:    GMProductItemId,
		})
		if err != nil {
			record.PrizeStatus = 2
			err2 := activity2.DefaultGMRecordRepository.Save(record)
			if err2 != nil {
				app.Logger.Error("返还兑换机会失败", userId, addressId, err, err2)
			}
			return nil, err
		}
		return order, nil
	}
	return nil, errors.New("状态错误,请联系管理员")
}
func (srv GMService) AnswerQuestion(param AnswerGMQuestionParam) (*activity.GMRecord, error) {
	record, err := srv.FindOrCreateGMRecord(param.UserId)
	if err != nil {
		return nil, err
	}
	if record.AvailableQuesNum <= 0 {
		return nil, errors.New("答题次数用光啦,快去邀请好友获取答题机会吧")
	}
	record.UsedQuesNum++
	if param.IsRight {
		record.RightQuesNum++
	} else {
		record.WrongQuesNum++
	}
	record.AvailableQuesNum--
	if record.RightQuesNum >= 5 {
		record.PrizeStatus = 2
	}
	err = activity2.DefaultGMRecordRepository.Save(record)
	if err != nil {
		app.Logger.Error("GM答题失败", param, err)
		return nil, errors.New("答题失败,请稍后再试")
	}
	isRight := 1
	if !param.IsRight {
		isRight = 2
	}

	quesLog := activity.GMQuestionLog{
		UserId:      param.UserId,
		Title:       param.Title,
		Answer:      param.Answer,
		IsRight:     isRight,
		IsSendPoint: 1,
		CreatedAt:   model.NewTime(),
		UpdatedAt:   model.NewTime(),
	}
	err = activity2.DefaultGMQuestionLogRepository.Save(&quesLog)
	if err != nil {
		app.Logger.Error("GM答题失败", param, err)
		return nil, errors.New("答题失败,请稍后再试")
	}

	//发放答题积分
	if param.IsRight {
		err = srv.SendAnswerQuestionBonus(param.UserId, quesLog.ID)
		if err != nil {
			return nil, err
		}
	}
	return record, nil
}

// SendAnswerQuestionBonus 发放积分
func (srv GMService) SendAnswerQuestionBonus(userId int64, logId int) error {
	user, err := service2.DefaultUserService.GetUserById(userId)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New("未查询到用户信息,请联系管理员")
	}

	log := activity2.DefaultGMQuestionLogRepository.FindById(logId)
	if log.ID == 0 {
		return errors.New("未查到答题记录")
	}
	if log.IsRight == 2 {
		return errors.New("答题错误,无法发放积分")
	}
	if log.IsSendPoint == 2 {
		return errors.New("积分已发放")
	}
	log.IsSendPoint = 2
	err = activity2.DefaultGMQuestionLogRepository.Save(&log)
	if err != nil {
		return err
	}

	_, err = service2.DefaultPointTransactionService.Create(service2.CreatePointTransactionParam{
		OpenId:       user.OpenId,
		Type:         entity2.POINT_QUIZ,
		Value:        50,
		AdditionInfo: fmt.Sprintf(`{"activity":"greenmonday","logId",%d}`, logId),
	})
	return err
}
func (srv GMService) AddInvitationRecord(userId, InviteeUserId int64) error {

	InviteeUser, err := service2.DefaultUserService.GetUserById(InviteeUserId)
	if err != nil {
		return err
	}
	if InviteeUser.ID == 0 {
		return errors.New("未查询到用户信息")
	}
	isNewUser := srv.IsNewUser(InviteeUser.Time.Time)

	record := activity2.DefaultGMInvitationRecordRepository.FindBy(activity2.FindGMInvitationRecordBy{
		InviteeUserId: InviteeUserId,
	})
	if record.ID == 0 {
		record = activity.GMInvitationRecord{
			UserId:           userId,
			InviteeUserId:    InviteeUserId,
			InviteeIsNewUser: util.TernaryOperator(isNewUser, 1, 2).Int(),
			CreatedAt:        model.NewTime(),
			UpdatedAt:        model.NewTime(),
		}
		err := activity2.DefaultGMInvitationRecordRepository.Create(&record)
		if err != nil {
			return errors.Wrap(err, "保存邀请记录失败")
		}

		if isNewUser {
			err := srv.AddQuestionNum(userId)
			if err != nil {
				app.Logger.Error("增加答题机会失败", err)
			}
		}
	}
	return nil
}
func (srv GMService) AddQuestionNum(userId int64) error {
	record, err := srv.FindOrCreateGMRecord(userId)
	if err != nil {
		return err
	}
	record.AvailableQuesNum++
	return activity2.DefaultGMRecordRepository.Save(record)
}
func (srv GMService) FindOrCreateGMRecord(userId int64) (*activity.GMRecord, error) {
	record := activity2.DefaultGMRecordRepository.FindBy(activity2.FindGMRecordBy{
		UserId: userId,
	})
	if record.ID != 0 {
		return &record, nil
	}

	record = activity.GMRecord{
		UserId:           userId,
		AvailableQuesNum: 1,
		UsedQuesNum:      0,
		RightQuesNum:     0,
		WrongQuesNum:     0,
		PrizeStatus:      1,
		CreatedAt:        model.NewTime(),
		UpdatedAt:        model.NewTime(),
	}
	err := activity2.DefaultGMRecordRepository.Create(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
func (srv GMService) IsNewUser(userCreatedTime time.Time) bool {
	return GMNewUserTime.Before(userCreatedTime)
}
