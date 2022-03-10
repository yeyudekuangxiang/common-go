package activity

import (
	"mio/model"
	activityM "mio/model/entity/activity"
	activityR "mio/repository/activity"
)

var DefaultBocShareBonusRecordService = BocShareBonusRecordService{}

type BocShareBonusRecordService struct {
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
	return &record, activityR.DefaultBocShareBonusRecordRepository.Save(&record)
}
