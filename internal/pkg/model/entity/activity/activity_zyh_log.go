package activity

import "mio/internal/pkg/model/entity"

type ZyhLog struct {
	Id             int64                       `json:"id"`
	Openid         string                      `json:"openid"`
	PointType      entity.PointTransactionType `json:"point_type"`
	Value          int64                       `json:"value"`
	ResultCode     string                      `json:"result_code"`
	AdditionalInfo string                      `json:"additional_info"`
	TransactionId  string                      `json:"transaction_id"`
}

func (ZyhLog) TableName() string {
	return "activity_zyh_log"
}
