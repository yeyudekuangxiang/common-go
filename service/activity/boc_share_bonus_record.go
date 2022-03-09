package activity

import (
	"mio/model"
	activityM "mio/model/entity/activity"
	activityR "mio/repository/activity"
)

var DefaultBocShareBonusRecordService = BocShareBonusRecordService{}

type BocShareBonusRecordService struct {
}

func (b BocShareBonusRecordService) CreateRecord(param CreateBocShareBonusRecordParam) (*activityM.BocShareBonusRecord, error) {
	record := activityM.BocShareBonusRecord{
		UserId:    param.UserId,
		Value:     param.Value,
		CreatedAt: model.NewTime(),
		UpdatedAt: model.NewTime(),
	}
	return &record, activityR.DefaultBocShareBonusRecordRepository.Save(&record)
}
