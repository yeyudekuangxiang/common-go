package srv_types

import "mio/internal/pkg/model/entity"

type GetZyhGetInfoByDTO struct {
	Openid string
	VolId  string
}

type GetZyhLogAddDTO struct {
	Openid         string                      `json:"openid"`
	PointType      entity.PointTransactionType `json:"point_type"`
	Value          int64                       `json:"value"`
	ResultCode     string                      `json:"result_code"`
	AdditionalInfo string                      `json:"additional_info"`
	TransactionId  string                      `json:"transaction_id"`
}
