package activity

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	entity2 "mio/internal/pkg/model/entity"
	activity2 "mio/internal/pkg/model/entity/activity"
	repository2 "mio/internal/pkg/repository"
	"mio/internal/pkg/repository/activity"
	service2 "mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"strconv"
	"time"
)

//活动开始时间 用户判断是否老用户
var bocActivityStartTime, _ = time.Parse("2006-01-02 15:04:05", "2022-03-09 00:00:00")

var DefaultBocService = BocService{repo: activity.DefaultBocRecordRepository}

type BocService struct {
	repo activity.BocRecordRepository
}

// GetApplyRecordPageList 获取用户的分享记录
func (b BocService) GetApplyRecordPageList(param GetRecordPageListParam) (list []BocRecordDetail, total int64, err error) {
	list = make([]BocRecordDetail, 0)
	//获取用户信息
	user, err := service2.DefaultUserService.GetUserById(param.UserId)
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
	userList := repository2.DefaultUserRepository.GetUserListBy(repository2.GetUserListBy{
		Mobile: user.PhoneNumber,
		Source: entity2.UserSourceMobile,
	})
	if len(userList) == 0 {
		return
	}

	//根据账户id列表查询所有的邀请记录
	userIds := make([]int64, 0)
	for _, u := range userList {
		userIds = append(userIds, u.ID)
	}
	recordList, total := activity.DefaultBocRecordRepository.GetPageListBy(activity.GetRecordListBy{
		ApplyStatus:             param.ApplyStatus,
		ShareUserIds:            userIds,
		ShareUserBocBonusStatus: param.ShareUserBocBonusStatus,
		Offset:                  param.Offset,
		Limit:                   param.Limit,
	})
	if len(recordList) == 0 {
		return
	}

	//根据记录中的用户id查询用户信息
	userIds = make([]int64, 0)
	for _, record := range recordList {
		userIds = append(userIds, record.UserId)
	}
	recordUserList := repository2.DefaultUserRepository.GetUserListBy(repository2.GetUserListBy{
		UserIds: userIds,
	})
	userMobiles := make([]string, 0)

	userMap := make(map[int64]entity2.User)
	for _, user := range recordUserList {
		userMap[user.ID] = user
		if user.PhoneNumber != "" {
			userMobiles = append(userMobiles, user.PhoneNumber)
		}
	}

	mioUserList := repository2.DefaultUserRepository.GetUserListBy(repository2.GetUserListBy{
		Mobiles: userMobiles,
		Source:  entity2.UserSourceMio,
	})
	mioUserMap := make(map[string]entity2.User)
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
			BocRecord:     record,
			UpdatedAtDate: record.UpdatedAt.Date(),
			CreatedAtDate: record.CreatedAt.Date(),
			User:          user.ShortUser(),
		})
	}
	return
}

// AddShareNum 添加用户的分享次数
func (b BocService) AddShareNum(userId int64) error {
	record := activity.DefaultBocRecordRepository.FindBy(activity.FindRecordBy{
		UserId: userId,
	})
	if record.Id == 0 {
		return nil
	}
	record.ShareNum++
	return activity.DefaultBocRecordRepository.Save(&record)
}

// FindApplyRecordMini 小程序查询用户参与的中行活动记录
func (b BocService) FindApplyRecordMini(userId int64) (*activity2.BocRecord, error) {
	defaultRecord := activity2.NewBocRecord()
	if userId == 0 {
		return &defaultRecord, nil
	}

	user, err := service2.DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 || user.PhoneNumber == "" {
		return &defaultRecord, nil
	}

	mobileUser := repository2.DefaultUserRepository.GetUserBy(repository2.GetUserBy{
		Mobile: user.PhoneNumber,
		Source: entity2.UserSourceMobile,
	})

	return b.FindOrCreateApplyRecord(AddApplyRecordParam{
		UserId: mobileUser.ID,
	})
}

// FindOrCreateApplyRecord 创建或者查询活动记录
func (b BocService) FindOrCreateApplyRecord(param AddApplyRecordParam) (*activity2.BocRecord, error) {
	if param.UserId == 0 {
		record := activity2.NewBocRecord()
		return &record, nil
	}

	record := activity.DefaultBocRecordRepository.FindBy(activity.FindRecordBy{
		UserId: param.UserId,
	})

	//已存在申请记录
	if record.Id > 0 {
		//奖励未发放
		if record.AnswerStatus == 2 && record.AnswerBonusStatus == 1 {
			go func() {
				if err := b.SendAnswerBonus(param.UserId); err != nil {
					app.Logger.Error("SendAnswerBonus", param.UserId, err)
				}
			}()
		}

		return &record, nil
	}

	record = activity2.BocRecord{
		UserId:                param.UserId,
		ShareUserId:           param.ShareUserId,
		ApplyStatus:           1,
		ApplyBonusStatus:      1,
		BindWechatStatus:      1,
		BindWechatBonusStatus: 1,
		AnswerStatus:          1,
		AnswerBonusStatus:     1,
		Source:                param.Source,
		CreatedAt:             model.Time{Time: time.Now()},
		UpdatedAt:             model.Time{Time: time.Now()},
	}

	if param.ShareUserId > 0 {
		_ = b.AddShareNum(param.ShareUserId)
	}

	return &record, activity.DefaultBocRecordRepository.Save(&record)
}

// SendAnswerBonus 发放用户积分
func (b BocService) SendAnswerBonus(userId int64) error {
	util.DefaultLock.LockWait(fmt.Sprintf("BocSendAnswerBonus%d", userId), 5*time.Second)
	defer func() {
		util.DefaultLock.UnLock(fmt.Sprintf("BocSendAnswerBonus%d", userId))
	}()

	if userId == 0 {
		return nil
	}

	record, err := b.FindOrCreateApplyRecord(AddApplyRecordParam{
		UserId: userId,
	})
	if err != nil {
		return err
	}
	if !(record.AnswerStatus == 2 && record.AnswerBonusStatus == 1) {
		return nil
	}

	mioUser, err := service2.DefaultUserService.FindUserBySource(entity2.UserSourceMio, userId)
	if err != nil {
		return err
	}

	if mioUser.ID == 0 {
		app.Logger.Warn("未绑定小程序,打开绿喵小程序绑定手机号后,打开小程序邀请记录页面将自动发放积分", userId)
		return nil
	}

	record.AnswerBonusStatus = 2
	record.UpdatedAt = model.NewTime()
	err = activity.DefaultBocRecordRepository.Save(record)
	if err != nil {
		return err
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

	_, err = service2.DefaultPointTransactionService.Create(service2.CreatePointTransactionParam{
		OpenId:       mioUser.OpenId,
		Type:         entity2.POINT_QUIZ,
		Value:        value,
		AdditionInfo: additionInfo,
	})
	return err
}

// AnswerQuestion 回答问题
func (b BocService) AnswerQuestion(userId int64, right int) error {
	util.DefaultLock.LockWait(fmt.Sprintf("BocAnswerQuestion%d", userId), 5*time.Second)
	defer func() {
		util.DefaultLock.UnLock(fmt.Sprintf("BocAnswerQuestion%d", userId))
	}()

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
	err = activity.DefaultBocRecordRepository.Save(record)
	if err != nil {
		return err
	}

	if record.AnswerStatus == 2 && record.AnswerBonusStatus == 1 {
		err := b.SendAnswerBonus(userId)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsOldUser 是否老用户
func (b BocService) IsOldUser(t time.Time) bool {
	return !t.IsZero() && t.Before(bocActivityStartTime)
}
func (b BocService) IsOldUserById(userId int64) (bool, error) {
	mioUser, err := service2.DefaultUserService.FindUserBySource(entity2.UserSourceMio, userId)
	if err != nil {
		return false, err
	}
	return b.IsOldUser(mioUser.Time.Time), nil
}

// SendApplyBonus 发放申请卡片奖励(拿到中行列表后 根据中行列表和活动参与记录执行)
func (b BocService) SendApplyBonus(userId int64) error {
	//防止并发
	cmd := app.Redis.GetEx(context.Background(), config.RedisKey.Limit1S+strconv.Itoa(int(userId)), 1*time.Second)
	fmt.Println(cmd)
	if cmd.Val() != "" {
		return errors.New("正在审核中,请稍等")
	}
	//更改奖励发放状态
	record := activity.DefaultBocRecordRepository.FindBy(activity.FindRecordBy{
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
	record.ApplyBonusTime = model.NewTime()
	record.UpdatedAt = model.NewTime()
	err := activity.DefaultBocRecordRepository.Save(&record)
	if err != nil {
		return err
	}

	//增加奖励发放记录
	user, err := service2.DefaultUserService.FindUserBySource(entity2.UserSourceMobile, userId)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return errors.New("请先绑定手机号后再进行操作")
	}

	_, err = DefaultBocShareBonusRecordService.CreateRecord(CreateBocShareBonusRecordParam{
		UserId: user.ID,
		Value:  500,
		Type:   activity2.BocShareBonusMio,
		Info:   "申请卡片获取5元话费奖励金",
	})
	if err != nil {
		return err
	}

	//后续进行实际话费充值操作
	userInfo, _ := service2.DefaultUserService.GetUserById(userId)                                                                 //需要手机号码
	err = service2.DefaultUnidianService.SendPrize(service2.UnidianTypeId.FiveYuan, userInfo.PhoneNumber, config.RedisKey.UniDian) //充话费
	if err != nil {
		return err
	}
	return nil
}

// SendBindWechatBonus 发放绑定微信奖励(拿到中行列表后 根据中行列表和活动参与记录执行)
func (b BocService) SendBindWechatBonus(userId int64) error {
	//更改奖励发放状态
	record := activity.DefaultBocRecordRepository.FindBy(activity.FindRecordBy{
		UserId: userId,
	})
	if record.Id == 0 {
		return errors.New("未查询到活动参与记录")
	}
	if record.ApplyBonusStatus == 2 {
		return errors.New("审核中")
	} else if record.ApplyBonusStatus == 3 {
		return errors.New("奖励已经发放过了")
	}
	record.BindWechatStatus = 2
	record.BindWechatBonusStatus = 2
	record.BindWechatBonusTime = model.NewTime()
	record.UpdatedAt = model.NewTime()
	err := activity.DefaultBocRecordRepository.Save(&record)
	if err != nil {
		return err
	}

	//增加奖励发放记录
	user, err := service2.DefaultUserService.FindUserBySource(entity2.UserSourceMobile, userId)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return errors.New("请先绑定手机号后再进行操作")
	}
	_, err = DefaultBocShareBonusRecordService.CreateRecord(CreateBocShareBonusRecordParam{
		UserId: user.ID,
		Value:  1000,
		Type:   activity2.BocShareBonusBoc10,
		Info:   "申请卡片获取10元消费金",
	})

	if err != nil {
		return err
	}

	//后续进行实际话费充值操作

	return nil
}

// ApplySendApplyBonus 申请发放五元奖励
func (b BocService) ApplySendApplyBonus(userId int64) error {
	mobileUser, err := service2.DefaultUserService.FindUserBySource(entity2.UserSourceMobile, userId)
	if err != nil {
		return err
	}
	if mobileUser.ID == 0 {
		return errors.New("请检查是否已经绑定手机号号码")
	}
	record := activity.DefaultBocRecordRepository.FindBy(activity.FindRecordBy{
		UserId: mobileUser.ID,
	})
	if record.ApplyStatus != 3 {
		return errors.New("银行卡暂未申请成功,请稍后再试")
	}
	if record.ApplyBonusStatus != 1 {
		return errors.New("奖励已申请")
	}
	record.ApplyBonusStatus = 2
	err = activity.DefaultBocRecordRepository.Save(&record)
	if err != nil {
		app.Logger.Error(userId, err)
		return errors.New("申请失败,请稍后再试")
	}

	//话费的无需审核直接充值
	return b.SendApplyBonus(userId)
}

// ApplySendBindWechatBonus 申请发放10圆奖励
func (b BocService) ApplySendBindWechatBonus(userId int64) error {
	mobileUser, err := service2.DefaultUserService.FindUserBySource(entity2.UserSourceMobile, userId)
	if err != nil {
		return err
	}
	if mobileUser.ID == 0 {
		return errors.New("请检查是否已经绑定手机号号码")
	}
	record := activity.DefaultBocRecordRepository.FindBy(activity.FindRecordBy{
		UserId: mobileUser.ID,
	})
	if record.BindWechatStatus != 2 {
		return errors.New("银行卡暂未绑定微信,请稍后再试")
	}
	if record.BindWechatBonusStatus != 1 {
		return errors.New("奖励已申请")
	}
	record.BindWechatBonusStatus = 2
	return activity.DefaultBocRecordRepository.Save(&record)
}

// ApplySendBocBonus 申请发放中行奖励金
func (b BocService) ApplySendBocBonus(userId int64) error {
	list, _, err := b.GetApplyRecordPageList(GetRecordPageListParam{
		UserId:                  userId,
		ApplyStatus:             3,
		ShareUserBocBonusStatus: 1,
		Offset:                  0,
		Limit:                   200,
	})
	if err != nil {
		app.Logger.Error(userId, err)
		return errors.New("申请失败,请稍后再试")
	}

	if len(list) == 0 {
		return errors.New("没有未领取的奖励")
	}

	sum, err := DefaultBocShareBonusRecordService.SendBocSum(userId)

	if err != nil {
		app.Logger.Error(userId, err)
		return errors.New("系统异常,请稍后再试")
	}
	if sum >= 5000 {
		app.Logger.Error(userId, sum, "奖励最高可领50元")
		return errors.New("奖励最高可领50元")
	}

	err = app.DB.Transaction(func(tx *gorm.DB) error {
		bocRecordRepo := activity.BocRecordRepository{DB: tx}
		shareBonusRecordService := BocShareBonusRecordService{repo: activity.BocShareBonusRecordRepository{DB: tx}}

		var totalValue int64
		ids := ""
		for _, item := range list {
			record := item.BocRecord
			record.ShareUserBocBonusStatus = 2
			record.ShareUserBocBonusTime = model.NewTime()
			err := bocRecordRepo.Save(&record)
			if err != nil {
				return err
			}
			ids += strconv.Itoa(int(record.Id)) + ","
			totalValue += 500

			if totalValue >= 5000-sum {
				totalValue = 5000 - sum
				break
			}
		}

		_, err = shareBonusRecordService.CreateRecord(CreateBocShareBonusRecordParam{
			UserId: userId,
			Value:  totalValue,
			Type:   activity2.BocShareBonusBoc,
			Info:   ids,
		})
		return err
	})
	if err != nil {
		app.Logger.Error(userId, err)
		return errors.New("申请失败,请稍后再试")
	}
	return nil
}
