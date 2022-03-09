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

//活动开始时间 用户判断是否老用户
var bocActivityStartTime, _ = time.Parse("2006-01-02 15:04:05", "2022-03-09 00:00:00")

var DefaultBocService = BocService{}

type BocService struct {
}

// GetApplyRecordPageList 获取用户的分享记录
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
	recordUserList := repository.DefaultUserRepository.GetUserListBy(repository.GetUserListBy{
		UserIds: userIds,
	})
	userMobiles := make([]string, 0)

	userMap := make(map[int64]entity.User)
	for _, user := range recordUserList {
		userMap[user.ID] = user
		if user.PhoneNumber != "" {
			userMobiles = append(userMobiles, user.PhoneNumber)
		}
	}

	mioUserList := repository.DefaultUserRepository.GetUserListBy(repository.GetUserListBy{
		Mobiles: userMobiles,
		Source:  entity.UserSourceMio,
	})
	mioUserMap := make(map[string]entity.User)
	for _, mioUser := range mioUserList {
		mioUserMap[mioUser.PhoneNumber] = mioUser
	}

	for _, record := range recordList {

		user := userMap[record.UserId]
		if mioUser, ok := mioUserMap[user.PhoneNumber]; ok {
			user.Nickname = mioUser.Nickname
			user.AvatarUrl = mioUser.AvatarUrl
		}

		list = append(list, BocRecordDetail{
			BocRecord: record,
			User:      user.ShortUser(),
		})
	}
	return
}

// AddShareNum 添加用户的分享次数
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

// FindApplyRecord 查询用户参与的中行活动记录
func (b BocService) FindApplyRecord(userId int64) (*activityM.BocRecord, error) {
	defaultRecord := activityM.NewBocRecord()
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

// FindOrCreateApplyRecord 创建或者查询活动记录
func (b BocService) FindOrCreateApplyRecord(param AddApplyRecordParam) (*activityM.BocRecord, error) {
	if param.UserId == 0 {
		record := activityM.NewBocRecord()
		return &record, nil
	}

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

// AnswerQuestion 回答问题
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

//回答问题领积分
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

	var (
		value        int
		additionInfo string
	)
	if b.IsOldUser(mioUser.Time.Time) {
		value = 500
		additionInfo = "老用户答题得500积分"
	} else {
		value = 2500
		additionInfo = "新用户答题得2500积分"
	}
	_, err = service.DefaultPointTransactionService.Create(service.CreatePointTransactionParam{
		OpenId:       mioUser.OpenId,
		Type:         entity.POINT_QUIZ,
		Value:        value,
		AdditionInfo: additionInfo,
	})
	return err
}

// IsOldUser 是否老用户
func (b BocService) IsOldUser(t time.Time) bool {
	return !t.IsZero() && t.Before(bocActivityStartTime)
}

// SendApplyBonus 发放申请卡片奖励
func (b BocService) SendApplyBonus(userId int64) error {
	//更改奖励发放状态
	record := activityR.DefaultBocRecordRepository.FindBy(activityR.FindRecordBy{
		UserId: userId,
	})
	if record.Id == 0 {
		return errors.New("未查询到活动参与记录")
	}

	if record.ApplyBonusStatus == 3 {
		return errors.New("奖励已经发放过了")
	}

	record.ApplyStatus = 3
	record.ApplyBonusStatus = 3
	record.UpdatedAt = model.NewTime()
	err := activityR.DefaultBocRecordRepository.Save(&record)
	if err != nil {
		return err
	}

	//增加奖励发放记录
	user, err := service.DefaultUserService.FindUserBySource(entity.UserSourceMobile, userId)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return errors.New("请先绑定手机号后再进行操作")
	}

	_, err = DefaultBocShareBonusRecordService.CreateRecord(CreateBocShareBonusRecordParam{
		UserId: user.ID,
		Value:  500,
		Type:   activityM.BocShareBonusMio,
		Info:   "申请卡片获取5元话费奖励金",
	})
	if err != nil {
		return err
	}

	//后续进行实际话费充值操作

	return nil
}

// SendBindWechatBonus 发放绑定微信奖励
func (b BocService) SendBindWechatBonus(userId int64) error {
	//更改奖励发放状态
	record := activityR.DefaultBocRecordRepository.FindBy(activityR.FindRecordBy{
		UserId: userId,
	})
	if record.Id == 0 {
		return errors.New("未查询到活动参与记录")
	}
	if record.ApplyBonusStatus == 3 {
		return errors.New("奖励已经发放过了")
	}
	record.BindWechatStatus = 2
	record.BindWechatBonusStatus = 3
	record.BindWechatBonusTime = model.NewTime()
	record.UpdatedAt = model.NewTime()
	err := activityR.DefaultBocRecordRepository.Save(&record)
	if err != nil {
		return err
	}

	//增加奖励发放记录
	user, err := service.DefaultUserService.FindUserBySource(entity.UserSourceMobile, userId)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return errors.New("请先绑定手机号后再进行操作")
	}
	_, err = DefaultBocShareBonusRecordService.CreateRecord(CreateBocShareBonusRecordParam{
		UserId: user.ID,
		Value:  1000,
		Type:   activityM.BocShareBonusMio,
		Info:   "申请卡片获取10元话费奖励金",
	})

	if err != nil {
		return err
	}

	//后续进行实际话费充值操作

	return nil
}
