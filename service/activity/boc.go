package activity

import (
	"github.com/pkg/errors"
	"mio/model/entity"
	"mio/repository"
	activityR "mio/repository/activity"
)

type BocService struct {
}

func (b BocService) GetRecordPageList(param GetRecordPageListParam) (list []BocApplyRecordDetail, total int64, err error) {
	list = make([]BocApplyRecordDetail, 0)
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

	userList := repository.DefaultUserRepository.GetUserListBy(repository.GetUserListBy{
		Mobile: user.PhoneNumber,
	})
	if len(userList) == 0 {
		return
	}

	userIds := make([]int64, 0)
	for _, u := range userList {
		userIds = append(userIds, u.ID)
	}
	recordList, total := activityR.DefaultBocApplyRecordRepository.GetPageListBy(activityR.GetRecordListBy{
		ApplyStatus: param.ApplyRecordStatus,
		UserIds:     userIds,
		Offset:      param.Offset,
		Limit:       param.Limit,
	})
	if len(recordList) == 0 {
		return
	}

	userIds = make([]int64, 0)
	for _, record := range recordList {
		userIds = append(userIds, record.UserId)
	}
	shareUserList := repository.DefaultUserRepository.GetUserListBy(repository.GetUserListBy{
		UserIds: userIds,
	})
	shareUserMap := make(map[int64]entity.User)
	for _, shareUser := range shareUserList {
		shareUserMap[shareUser.ID] = shareUser
	}

	for _, record := range recordList {
		list = append(list, BocApplyRecordDetail{
			BocApplyRecord: record,
			ShareUser:      shareUserMap[record.UserId],
		})
	}
	return
}
