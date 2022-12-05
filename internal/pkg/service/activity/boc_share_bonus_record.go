package activity

import (
	"mio/internal/pkg/model"
	activityM "mio/internal/pkg/model/entity/activity"
	activityR "mio/internal/pkg/repository/activity"
)

var DefaultBocShareBonusRecordService = BocShareBonusRecordService{repo: activityR.DefaultBocShareBonusRecordRepository}

type BocShareBonusRecordService struct {
	repo activityR.BocShareBonusRecordRepository
}

// CreateRecord 添加中行活动奖励领取记录 申请卡片和绑定微信话费奖励发放后 需要将boc_record表的相关信息
func (b BocShareBonusRecordService) CreateRecord(param CreateBocShareBonusRecordParam) (*activityM.BocShareBonusRecord, error) {
	record := activityM.BocShareBonusRecord{
		UserId:    param.UserId,
		Value:     param.Value,
		Type:      param.Type,
		Info:      param.Info,
		CreatedAt: model.NewTime(),
		UpdatedAt: model.NewTime(),
	}

	return &record, b.repo.Save(&record)
}
func (b BocShareBonusRecordService) SendBocSum(userId int64) (int64, error) {
	return b.repo.GetUserBonus(userId, activityM.BocShareBonusBoc), nil
}