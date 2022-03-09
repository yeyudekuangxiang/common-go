package activity

import (
	"github.com/pkg/errors"
	"mio/model"
	"mio/model/entity"
	activityM "mio/model/entity/activity"
	"mio/repository"
	activityR "mio/repository/activity"
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
	if user.ID == 0 {
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
