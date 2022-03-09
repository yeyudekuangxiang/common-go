package activity

import (
	"github.com/pkg/errors"
	"mio/model"
	"mio/model/entity"
	activityM "mio/model/entity/activity"
	"mio/repository"
	activityR "mio/repository/activity"
	"mio/service"
	"time"
)

var DefaultBocService = BocService{}

type BocService struct {
}

func (b BocService) GetApplyRecordPageList(param GetRecordPageListParam) (list []BocRecordDetail, total int64, err error) {
	list = make([]BocRecordDetail, 0)
	//获取用户信息
	user, err := repository.DefaultUserRepository.GetUserById(param.UserId)
	if err != nil {
		return
	}
	if user == nil {
		err = errors.New("用户不存在")
		return
	}
	if user.PhoneNumber == "" {
		err = errors.New("未绑定手机号,请先绑定手机号")
		return
	}

	//根据手机号获取用户所有账户
	userList := repository.DefaultUserRepository.GetUserListBy(repository.GetUserListBy{
		Mobile: user.PhoneNumber,
	})
	if len(userList) == 0 {
		return
	}

	//根据账户id列表查询所有的邀请记录
	userIds := make([]int64, 0)
	for _, u := range userList {
		userIds = append(userIds, u.ID)
	}
	recordList, total := activityR.DefaultBocRecordRepository.GetPageListBy(activityR.GetRecordListBy{
		ApplyStatus:  param.ApplyStatus,
		ShareUserIds: userIds,
		Offset:       param.Offset,
		Limit:        param.Limit,
	})
	if len(recordList) == 0 {
		return
	}

	//根据记录中的用户id查询用户信息
	userIds = make([]int64, 0)
	for _, record := range recordList {
		userIds = append(userIds, record.UserId)
	}
	recordUserList := repository.DefaultUserRepository.GetShortUserListBy(repository.GetUserListBy{
		UserIds: userIds,
	})
	userMap := make(map[int64]entity.ShortUser)
	for _, user := range recordUserList {
		userMap[user.ID] = user
	}

	for _, record := range recordList {
		list = append(list, BocRecordDetail{
			BocRecord: record,
			User:      userMap[record.UserId],
		})
	}
	return
}
func (b BocService) AddShareNum(userId int64) error {
	record := activityR.DefaultBocRecordRepository.FindBy(activityR.FindRecordBy{
		UserId: userId,
	})
	if record.Id == 0 {
		return nil
	}
	record.ShareNum++
	return activityR.DefaultBocRecordRepository.Save(&record)
}
func (b BocService) AddApplyRecord(param AddApplyRecordParam) (*activityM.BocRecord, error) {
	record := activityR.DefaultBocRecordRepository.FindBy(activityR.FindRecordBy{
		UserId: param.UserId,
	})
	//已存在申请记录
	if record.Id > 0 {
		return &record, nil
	}

	record = activityM.BocRecord{
		UserId:                param.UserId,
		ShareUserId:           param.ShareUserId,
		ApplyStatus:           1,
		ApplyBonusStatus:      1,
		BindWechatStatus:      1,
		BindWechatBonusStatus: 1,
		AnswerStatus:          1,
		AnswerBonusStatus:     1,
		CreatedAt:             model.Time{Time: time.Now()},
		UpdatedAt:             model.Time{Time: time.Now()},
	}

	if param.ShareUserId > 0 {
		_ = b.AddShareNum(param.ShareUserId)
	}

	return &record, activityR.DefaultBocRecordRepository.Save(&record)
}
func (b BocService) FindApplyRecord(userId int64) (*activityM.BocRecord, error) {

	defaultRecord := activityM.NewBocRecord()
	defaultRecord.UserId = userId
	if userId == 0 {
		return &defaultRecord, nil
	}

	user, err := repository.DefaultUserRepository.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	if user == nil || user.PhoneNumber == "" {
		return &defaultRecord, nil
	}

	mobileUser := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{
		Mobile: user.PhoneNumber,
		Source: entity.UserSourceMobile,
	})

	return b.FindOrCreateApplyRecord(AddApplyRecordParam{
		UserId: mobileUser.ID,
	})
}
func (b BocService) FindOrCreateApplyRecord(param AddApplyRecordParam) (*activityM.BocRecord, error) {
	if param.UserId == 0 {
		record := activityM.NewBocRecord()
		return &record, nil
	}

	record := activityR.DefaultBocRecordRepository.FindBy(activityR.FindRecordBy{
		UserId: param.UserId,
	})

	if record.Id == 0 {
		return b.AddApplyRecord(param)
	}
	return &record, nil
}
func (b BocService) AnswerQuestion(userId int64, right int) error {
	record, err := b.FindOrCreateApplyRecord(AddApplyRecordParam{
		UserId: userId,
	})
	if err != nil {
		return err
	}

	if record.AnswerStatus > 1 {
		return nil
	}

	record.AnswerStatus = right
	record.UpdatedAt = model.NewTime()
	record.AnswerBonusStatus = 2
	err = activityR.DefaultBocRecordRepository.Save(record)
	if err != nil {
		return err
	}
	if right == 2 {
		return b.makeAnswerPointTransaction(userId)
	}
	return nil
}
func (b BocService) makeAnswerPointTransaction(userId int64) error {
	user, err := repository.DefaultUserRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	if user == nil || user.PhoneNumber == "" {
		return errors.New("积分发放失败,用户不存在")
	}

	mioUser := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{
		Mobile: user.PhoneNumber,
		Source: entity.UserSourceMio,
	})

	if mioUser.ID == 0 {
		return errors.New("积分发放失败,小程序未绑定手机号,请在小程序端绑定手机号")
	}

	_, err = service.DefaultPointTransactionService.Create(service.CreatePointTransactionParam{
		OpenId:       mioUser.OpenId,
		Type:         entity.POINT_QUIZ,
		Value:        2500,
		AdditionInfo: "",
	})
	return err
}
