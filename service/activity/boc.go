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

type BocService struct {
}

func (b BocService) GetApplyRecordPageList(param GetRecordPageListParam) (list []BocApplyRecordDetail, total int64, err error) {
	list = make([]BocApplyRecordDetail, 0)
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
	recordList, total := activityR.DefaultBocApplyRecordRepository.GetPageListBy(activityR.GetRecordListBy{
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
		list = append(list, BocApplyRecordDetail{
			BocApplyRecord: record,
			User:           userMap[record.UserId],
		})
	}
	return
}
func (b BocService) AddApplyRecord(param AddApplyRecordParam) (*activityM.BocApplyRecord, error) {
	record := activityR.DefaultBocApplyRecordRepository.FindBy(activityR.FindRecordBy{
		UserId: param.UserId,
	})
	//已存在申请记录
	if record.Id > 0 {
		return &record, nil
	}

	record = activityM.BocApplyRecord{
		UserId:           param.UserId,
		ShareUserId:      param.ShareUserId,
		ApplyStatus:      1,
		IsGiveApplyBonus: 1,
		IsBindWechat:     1,
		IsGiveBindBonus:  1,
		CreatedAt:        model.Time{Time: time.Now()},
		UpdatedAt:        model.Time{Time: time.Now()},
	}
	return &record, activityR.DefaultBocApplyRecordRepository.Save(&record)
}
